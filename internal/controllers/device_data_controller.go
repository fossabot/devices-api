package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	dagrpc "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DIMO-Network/shared"

	"github.com/DIMO-Network/devices-api/internal/services"
	"github.com/segmentio/ksuid"

	"github.com/DIMO-Network/devices-api/internal/appmetrics"
	"github.com/DIMO-Network/devices-api/internal/constants"
	"github.com/DIMO-Network/devices-api/internal/controllers/helpers"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	smartcar "github.com/smartcar/go-sdk"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type QueryDeviceErrorCodesReq struct {
	ErrorCodes []string `json:"errorCodes" example:"P0106,P0279"`
}

type QueryDeviceErrorCodesResponse struct {
	ErrorCodes []services.ErrorCodesResponse `json:"errorCodes"`
	ClearedAt  *time.Time                    `json:"clearedAt"`
}

type GetUserDeviceErrorCodeQueriesResponse struct {
	Queries []GetUserDeviceErrorCodeQueriesResponseItem `json:"queries"`
}

type GetUserDeviceErrorCodeQueriesResponseItem struct {
	ErrorCodes  []services.ErrorCodesResponse `json:"errorCodes"`
	RequestedAt time.Time                     `json:"requestedAt" example:"2023-05-23T12:56:36Z"`
	// ClearedAt is the time at which the user cleared the codes from this query.
	// May be null.
	ClearedAt *time.Time `json:"clearedAt" example:"2023-05-23T12:57:05Z"`
}

// calculateRange returns the current estimated range based on fuel tank capacity, mpg, and fuelPercentRemaining and returns it in Kilometers
func calculateRange(ctx context.Context, ddSvc services.DeviceDefinitionService, deviceDefinitionID string, deviceStyleID null.String, fuelPercentRemaining float64) (*float64, error) {
	if fuelPercentRemaining <= 0.01 {
		return nil, errors.New("fuelPercentRemaining lt 0.01 so cannot calculate range")
	}

	dd, err := ddSvc.GetDeviceDefinitionByID(ctx, deviceDefinitionID)

	if err != nil {
		return nil, helpers.GrpcErrorToFiber(err, "deviceDefSvc error getting definition id: "+deviceDefinitionID)
	}

	rangeData := helpers.GetActualDeviceDefinitionMetadataValues(dd, deviceStyleID)

	// calculate, convert to Km
	if rangeData.FuelTankCapGal > 0 && rangeData.Mpg > 0 {
		fuelTankAtGal := rangeData.FuelTankCapGal * fuelPercentRemaining
		rangeMiles := rangeData.Mpg * fuelTankAtGal
		rangeKm := 1.60934 * rangeMiles
		return &rangeKm, nil
	}

	return nil, nil
}

// GetUserDeviceStatus godoc
// @Description Returns the latest status update for the device. May return 404 if the
// @Description user does not have a device with the ID, or if no status updates have come. Note this endpoint also exists under nft_controllers
// @Tags        user-devices
// @Produce     json
// @Param       user_device_id path     string true "user device ID"
// @Success     200            {object} controllers.DeviceSnapshot
// @Security    BearerAuth
// @Router      /user/devices/{userDeviceID}/status [get]
func (udc *UserDevicesController) GetUserDeviceStatus(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")
	userID := helpers.GetUserID(c)

	userDevice, err := models.FindUserDevice(c.Context(), udc.DBS().Reader, userDeviceID)
	if err != nil {
		return err
	}
	if userDevice.UserID != userID {
		return fiber.NewError(fiber.StatusForbidden)
	}

	udd, err := udc.deviceDataSvc.GetDeviceData(c.Context(),
		userDeviceID,
		userDevice.DeviceDefinitionID,
		userDevice.DeviceStyleID.String,
		[]int64{NonLocationData, CurrentLocation, AllTimeLocation}, // assume all privileges when called from here
	)
	if err != nil {
		err := shared.GrpcErrorToFiber(err, "failed to get user device data grpc")
		if err, ok := err.(*fiber.Error); ok && err.Code == 404 {
			helpers.SkipErrorLog(c)
		}
		return err
	}

	ds := grpcDeviceDataToSnapshot(udd)

	return c.JSON(ds)
}

func grpcDeviceDataToSnapshot(udd *dagrpc.UserDeviceDataResponse) DeviceSnapshot {
	ds := DeviceSnapshot{
		Charging:             udd.Charging,
		FuelPercentRemaining: udd.FuelPercentRemaining,
		BatteryCapacity:      udd.BatteryCapacity,
		OilLevel:             udd.OilLevel,
		Odometer:             udd.Odometer,
		Latitude:             udd.Latitude,
		Longitude:            udd.Longitude,
		Range:                udd.Range,
		StateOfCharge:        udd.StateOfCharge,
		ChargeLimit:          udd.ChargeLimit,
		RecordUpdatedAt:      convertTimestamp(udd.RecordUpdatedAt),
		RecordCreatedAt:      convertTimestamp(udd.RecordCreatedAt),
		TirePressure:         convertTirePressure(udd.TirePressure),
		BatteryVoltage:       udd.BatteryVoltage,
		AmbientTemp:          udd.AmbientTemp,
	}
	return ds
}

func convertTimestamp(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func convertTirePressure(tp *dagrpc.TirePressureResponse) *smartcar.TirePressure {
	if tp == nil {
		return nil
	}
	return &smartcar.TirePressure{
		FrontLeft:  tp.FrontLeft,
		FrontRight: tp.FrontRight,
		BackLeft:   tp.BackLeft,
		BackRight:  tp.BackRight,
	}
}

// RefreshUserDeviceStatus godoc
// @Description Starts the process of refreshing device status from Smartcar
// @Tags        user-devices
// @Param       user_device_id path string true "user device ID"
// @Success     204
// @Failure     429 "rate limit hit for integration"
// @Security    BearerAuth
// @Router      /user/devices/{userDeviceID}/commands/refresh [post]
func (udc *UserDevicesController) RefreshUserDeviceStatus(c *fiber.Ctx) error {
	udi := c.Params("userDeviceID")
	// We could probably do a smarter join here, but it's unclear to me how to handle that
	// in SQLBoiler.
	ud, err := models.UserDevices(
		models.UserDeviceWhere.ID.EQ(udi),
		qm.Load(models.UserDeviceRels.UserDeviceData),
	).One(c.Context(), udc.DBS().Reader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return err
	}
	smartCarInteg, err := udc.DeviceDefSvc.GetIntegrationByVendor(c.Context(), constants.SmartCarVendor)
	if err != nil {
		return helpers.GrpcErrorToFiber(err, "failed to get smartcar integration")
	}

	for _, deviceDatum := range ud.R.UserDeviceData {
		if deviceDatum.IntegrationID == smartCarInteg.Id {
			nextAvailableTime := deviceDatum.UpdatedAt.Add(time.Second * time.Duration(smartCarInteg.RefreshLimitSecs))
			if time.Now().Before(nextAvailableTime) {
				return fiber.NewError(fiber.StatusTooManyRequests, "rate limit for integration refresh hit")
			}

			udai, err := models.FindUserDeviceAPIIntegration(c.Context(), udc.DBS().Reader, deviceDatum.UserDeviceID, deviceDatum.IntegrationID)
			if err != nil {
				return err
			}
			if udai.Status == models.UserDeviceAPIIntegrationStatusActive && udai.TaskID.Valid {
				err = udc.smartcarTaskSvc.Refresh(udai)
				if err != nil {
					return err
				}
				return c.SendStatus(204)
			}

			return fiber.NewError(fiber.StatusConflict, "Integration not active.")
		}
	}
	return fiber.NewError(fiber.StatusBadRequest, "no active Smartcar integration found for this device")
}

var errorCodeRegex = regexp.MustCompile(`^.{5,8}$`)

// QueryDeviceErrorCodes godoc
// @Summary     Obtain, store, and return descriptions for a list of error codes from this vehicle.
// @Tags        error-codes
// @Param       userDeviceID path string true "user device id"
// @Param       queryDeviceErrorCodes body controllers.QueryDeviceErrorCodesReq true "error codes"
// @Success     200 {object} controllers.QueryDeviceErrorCodesResponse
// @Failure     404 {object} helpers.ErrorRes "Vehicle not found"
// @Security    BearerAuth
// @Router      /user/devices/{userDeviceID}/error-codes [post]
func (udc *UserDevicesController) QueryDeviceErrorCodes(c *fiber.Ctx) error {
	udi := c.Params("userDeviceID")

	logger := helpers.GetLogger(c, udc.log)
	ud, err := models.UserDevices(models.UserDeviceWhere.ID.EQ(udi)).One(c.Context(), udc.DBS().Reader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "No device with that id found.")
		}
		return err
	}

	dd, err := udc.DeviceDefSvc.GetDeviceDefinitionByID(c.Context(), ud.DeviceDefinitionID)
	if err != nil {
		return helpers.GrpcErrorToFiber(err, "deviceDefSvc error getting definition id: "+ud.DeviceDefinitionID)
	}

	req := &QueryDeviceErrorCodesReq{}
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Couldn't parse request.")
	}

	errorCodesLimit := 100
	if len(req.ErrorCodes) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "No error codes provided")
	}
	if len(req.ErrorCodes) > errorCodesLimit {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Too many error codes. Error codes list must be %d or below in length.", errorCodesLimit))
	}

	for _, v := range req.ErrorCodes {
		if !errorCodeRegex.MatchString(v) {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid error code %s", v))
		}
	}

	appmetrics.OpenAITotalCallsOps.Inc() // record new total call to chatgpt
	chtResp, err := udc.openAI.GetErrorCodesDescription(dd.Type.Make, dd.Type.Model, req.ErrorCodes)
	if err != nil {
		appmetrics.OpenAITotalFailedCallsOps.Inc()
		logger.Err(err).Interface("requestBody", req).Msg("Error occurred fetching description for error codes")
		return err
	}

	chtJSON, err := json.Marshal(chtResp)
	if err != nil {
		logger.Err(err).Interface("requestBody", req).Msg("Error occurred fetching description for error codes")
		return fiber.NewError(fiber.StatusInternalServerError, "Error occurred fetching description for error codes")
	}

	q := &models.ErrorCodeQuery{ID: ksuid.New().String(), UserDeviceID: udi, CodesQueryResponse: null.JSONFrom(chtJSON)}
	err = q.Insert(c.Context(), udc.DBS().Writer, boil.Infer())

	if err != nil {
		// TODO - should we return an error for this or just log it
		logger.Err(err).Msg("Could not save user query response")
	}

	return c.JSON(&QueryDeviceErrorCodesResponse{
		ErrorCodes: chtResp,
	})
}

// GetUserDeviceErrorCodeQueries godoc
// @Summary List all error code queries made for this vehicle.
// @Tags        error-codes
// @Param       userDeviceID path string true "user device id"
// @Success     200 {object} controllers.GetUserDeviceErrorCodeQueriesResponse
// @Failure     404 {object} helpers.ErrorRes "Vehicle not found"
// @Security    BearerAuth
// @Router      /user/devices/{userDeviceID}/error-codes [get]
func (udc *UserDevicesController) GetUserDeviceErrorCodeQueries(c *fiber.Ctx) error {
	logger := helpers.GetLogger(c, udc.log)

	userDeviceID := c.Params("userDeviceID")

	userDevice, err := models.UserDevices(
		models.UserDeviceWhere.ID.EQ(userDeviceID),
		qm.Load(models.UserDeviceRels.ErrorCodeQueries, qm.OrderBy(models.ErrorCodeQueryColumns.CreatedAt+" DESC")),
	).One(c.Context(), udc.DBS().Reader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "Could not find user device")
		}
		logger.Err(err).Msg("error occurred when fetching error codes for device")
		return fiber.NewError(fiber.StatusInternalServerError, "error occurred fetching device error queries")
	}

	queries := []GetUserDeviceErrorCodeQueriesResponseItem{}

	for _, erc := range userDevice.R.ErrorCodeQueries {
		ercJSON := []services.ErrorCodesResponse{}
		if err := erc.CodesQueryResponse.Unmarshal(&ercJSON); err != nil {
			return err
		}

		userDeviceresp := GetUserDeviceErrorCodeQueriesResponseItem{
			ErrorCodes:  ercJSON,
			RequestedAt: erc.CreatedAt,
			ClearedAt:   erc.ClearedAt.Ptr(),
		}

		queries = append(queries, userDeviceresp)
	}

	return c.JSON(GetUserDeviceErrorCodeQueriesResponse{Queries: queries})
}

// ClearUserDeviceErrorCodeQuery godoc
// @Summary     Mark the most recent set of error codes as having been cleared.
// @Tags        error-codes
// @Success     200 {object} controllers.QueryDeviceErrorCodesResponse
// @Failure     429 {object} helpers.ErrorRes "Last query already cleared"
// @Failure     404 {object} helpers.ErrorRes "Vehicle not found"
// @Security    BearerAuth
// @Router      /user/devices/{userDeviceID}/error-codes/clear [post]
func (udc *UserDevicesController) ClearUserDeviceErrorCodeQuery(c *fiber.Ctx) error {
	udi := c.Params("userDeviceID")

	logger := helpers.GetLogger(c, udc.log)

	errCodeQuery, err := models.ErrorCodeQueries(
		models.ErrorCodeQueryWhere.UserDeviceID.EQ(udi),
		qm.OrderBy(models.ErrorCodeQueryColumns.CreatedAt+" DESC"),
		qm.Limit(1),
	).One(c.Context(), udc.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("error occurred when fetching error codes for device")
		return fiber.NewError(fiber.StatusBadRequest, "error occurred fetching device error queries")
	}

	if errCodeQuery.ClearedAt.Valid {
		return fiber.NewError(fiber.StatusBadRequest, "all error codes already cleared")
	}

	errCodeQuery.ClearedAt = null.TimeFrom(time.Now().UTC().Truncate(time.Microsecond))
	if _, err = errCodeQuery.Update(c.Context(), udc.DBS().Writer, boil.Whitelist(models.ErrorCodeQueryColumns.ClearedAt)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error occurred updating device error queries")
	}

	errorCodeResp := []services.ErrorCodesResponse{}
	if err := errCodeQuery.CodesQueryResponse.Unmarshal(&errorCodeResp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error occurred updating device error queries")
	}

	return c.JSON(&QueryDeviceErrorCodesResponse{
		ErrorCodes: errorCodeResp,
		ClearedAt:  &errCodeQuery.ClearedAt.Time,
	})
}

// DeviceSnapshot is the response object for device status endpoint
// https://docs.google.com/document/d/1DYzzTOR9WA6WJNoBnwpKOoxfmrVwPWNLv0x0MkjIAqY/edit#heading=h.dnp7xngl47bw
type DeviceSnapshot struct {
	Charging             *bool                  `json:"charging,omitempty"`
	FuelPercentRemaining *float64               `json:"fuelPercentRemaining,omitempty"`
	BatteryCapacity      *int64                 `json:"batteryCapacity,omitempty"`
	OilLevel             *float64               `json:"oil,omitempty"`
	Odometer             *float64               `json:"odometer,omitempty"`
	Latitude             *float64               `json:"latitude,omitempty"`
	Longitude            *float64               `json:"longitude,omitempty"`
	Range                *float64               `json:"range,omitempty"`
	StateOfCharge        *float64               `json:"soc,omitempty"`
	ChargeLimit          *float64               `json:"chargeLimit,omitempty"`
	RecordUpdatedAt      *time.Time             `json:"recordUpdatedAt,omitempty"`
	RecordCreatedAt      *time.Time             `json:"recordCreatedAt,omitempty"`
	TirePressure         *smartcar.TirePressure `json:"tirePressure,omitempty"`
	BatteryVoltage       *float64               `json:"batteryVoltage,omitempty"`
	AmbientTemp          *float64               `json:"ambientTemp,omitempty"`
}
