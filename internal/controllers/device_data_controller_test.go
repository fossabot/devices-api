package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	dagrpc "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"github.com/tidwall/gjson"

	"github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/devices-api/internal/config"
	"github.com/DIMO-Network/devices-api/internal/constants"
	"github.com/DIMO-Network/devices-api/internal/services"
	mock_services "github.com/DIMO-Network/devices-api/internal/services/mocks"
	"github.com/DIMO-Network/devices-api/internal/test"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const migrationsDirRelPath = "../../migrations"

func TestUserDevicesController_calculateRange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	deviceDefSvc := mock_services.NewMockDeviceDefinitionService(mockCtrl)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

	ddID := ksuid.New().String()
	styleID := null.StringFrom(ksuid.New().String())
	attrs := []*grpc.DeviceTypeAttribute{
		{
			Name:  "fuel_tank_capacity_gal",
			Value: "15",
		},
		{
			Name:  "mpg",
			Value: "20",
		},
	}
	deviceDefSvc.EXPECT().GetDeviceDefinitionByID(gomock.Any(), ddID).Times(1).Return(&grpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: ddID,
		Verified:           true,
		DeviceAttributes:   attrs,
	}, nil)

	_ = NewUserDevicesController(&config.Settings{Port: "3000"}, nil, &logger, deviceDefSvc, nil, &fakeEventService{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	rge, err := calculateRange(ctx, deviceDefSvc, ddID, styleID, .7)
	require.NoError(t, err)
	require.NotNil(t, rge)
	assert.Equal(t, 337.9614, *rge)
}

type deps struct {
	deviceDefIntSvc        *mock_services.MockDeviceDefinitionIntegrationService
	deviceDefSvc           *mock_services.MockDeviceDefinitionService
	scClient               *mock_services.MockSmartcarClient
	scTaskSvc              *mock_services.MockSmartcarTaskService
	teslaSvc               *mock_services.MockTeslaService
	teslaTaskService       *mock_services.MockTeslaTaskService
	nhtsaService           *mock_services.MockINHTSAService
	autoPiIngest           *mock_services.MockIngestRegistrar
	deviceDefinitionIngest *mock_services.MockDeviceDefinitionRegistrar
	autoPiTaskSvc          *mock_services.MockAutoPiTaskService
	openAISvc              *mock_services.MockOpenAI
	logger                 zerolog.Logger
	mockCtrl               *gomock.Controller
	credentialSvc          *mock_services.MockVCService
	deviceDataSvc          *mock_services.MockDeviceDataService
}

func createMockDependencies(t *testing.T) deps {
	// arrange global db and route setup
	mockCtrl := gomock.NewController(t)

	deviceDefIntSvc := mock_services.NewMockDeviceDefinitionIntegrationService(mockCtrl)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionService(mockCtrl)
	deviceDataSvc := mock_services.NewMockDeviceDataService(mockCtrl)
	scClient := mock_services.NewMockSmartcarClient(mockCtrl)
	scTaskSvc := mock_services.NewMockSmartcarTaskService(mockCtrl)
	teslaSvc := mock_services.NewMockTeslaService(mockCtrl)
	teslaTaskService := mock_services.NewMockTeslaTaskService(mockCtrl)
	nhtsaService := mock_services.NewMockINHTSAService(mockCtrl)
	autoPiIngest := mock_services.NewMockIngestRegistrar(mockCtrl)
	deviceDefinitionIngest := mock_services.NewMockDeviceDefinitionRegistrar(mockCtrl)
	autoPiTaskSvc := mock_services.NewMockAutoPiTaskService(mockCtrl)
	openAISvc := mock_services.NewMockOpenAI(mockCtrl)
	credentialSvc := mock_services.NewMockVCService(mockCtrl)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

	return deps{
		deviceDefIntSvc:        deviceDefIntSvc,
		deviceDefSvc:           deviceDefSvc,
		scClient:               scClient,
		scTaskSvc:              scTaskSvc,
		teslaSvc:               teslaSvc,
		teslaTaskService:       teslaTaskService,
		nhtsaService:           nhtsaService,
		autoPiIngest:           autoPiIngest,
		deviceDefinitionIngest: deviceDefinitionIngest,
		autoPiTaskSvc:          autoPiTaskSvc,
		openAISvc:              openAISvc,
		logger:                 logger,
		mockCtrl:               mockCtrl,
		credentialSvc:          credentialSvc,
		deviceDataSvc:          deviceDataSvc,
	}

}

// Device Data Tests
func TestUserDevicesController_GetUserDeviceStatus(t *testing.T) {
	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, mockDeps.deviceDataSvc, nil)
	app := fiber.New()
	app.Get("/user/devices/:userDeviceID/status", test.AuthInjectorTestHandler(testUserID), c.GetUserDeviceStatus)

	t.Run("GET - get device status happy path", func(t *testing.T) {
		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)
		fuel := 0.50
		odo := 3000.50
		volt := 13.3
		mockDeps.deviceDataSvc.EXPECT().GetDeviceData(gomock.Any(), ud.ID, ud.DeviceDefinitionID,
			ud.DeviceStyleID.String, []int64{1, 3, 4}).Times(1).
			Return(&dagrpc.UserDeviceDataResponse{
				FuelPercentRemaining: &fuel,
				Odometer:             &odo,
				RecordUpdatedAt:      nil,
				RecordCreatedAt:      nil,
				BatteryVoltage:       &volt,
			}, nil)

		request := test.BuildRequest("GET", "/user/devices/"+ud.ID+"/status", "")
		response, err := app.Test(request, 60*1000)
		require.NoError(t, err, "failed to make request")
		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		body, _ := io.ReadAll(response.Body)
		if response.StatusCode != fiber.StatusOK {
			fmt.Println("body response: " + string(body))
		}

		assert.Equal(t, 3000.50, gjson.GetBytes(body, "odometer").Float())
		assert.Equal(t, 0.50, gjson.GetBytes(body, "fuelPercentRemaining").Float())

	})

}

// QueryErrorCodes Test
func TestUserDevicesController_QueryDeviceErrorCodes(t *testing.T) {

	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.QueryDeviceErrorCodes)

	t.Run("POST - get description for query codes", func(t *testing.T) {
		req := QueryDeviceErrorCodesReq{
			ErrorCodes: []string{"P0017", "P0016"},
		}

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		mockDeps.deviceDefSvc.
			EXPECT().
			GetDeviceDefinitionByID(gomock.Any(), ud.DeviceDefinitionID).
			Return(&grpc.GetDeviceDefinitionItemResponse{
				Type: &grpc.DeviceType{
					Make:  "Toyota",
					Model: "Camry",
					Year:  2023,
				},
			}, nil).
			AnyTimes()

		openAIResp := []services.ErrorCodesResponse{
			{
				Code:        "P0113",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}

		mockDeps.openAISvc.
			EXPECT().
			GetErrorCodesDescription(gomock.Eq("Toyota"), gomock.Eq("Camry"), gomock.Eq(req.ErrorCodes)).
			Return(openAIResp, nil).
			AnyTimes()

		j, _ := json.Marshal(req)

		request := test.BuildRequest("POST", "/user/devices/"+ud.ID+"/error-codes", string(j))
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		chatGptResp := QueryDeviceErrorCodesResponse{
			ErrorCodes: openAIResp,
		}
		chtJSON, err := json.Marshal(chatGptResp)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		assert.Equal(t,
			chtJSON,
			body,
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ShouldErrorOnTooManyErrorCodes(t *testing.T) {
	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.QueryDeviceErrorCodes)

	t.Run("POST - get description for query codes", func(t *testing.T) {

		erCodes := []string{}
		for i := 10; i <= 120; i++ {
			erCodes = append(erCodes, fmt.Sprintf("P000%d", i))
		}
		req := QueryDeviceErrorCodesReq{
			ErrorCodes: erCodes,
		}

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		mockDeps.deviceDefSvc.
			EXPECT().
			GetDeviceDefinitionByID(gomock.Any(), ud.DeviceDefinitionID).
			Return(&grpc.GetDeviceDefinitionItemResponse{
				Type: &grpc.DeviceType{
					Make:  "Toyota",
					Model: "Camry",
					Year:  2023,
				},
			}, nil).
			AnyTimes()

		chatGptResp := []services.ErrorCodesResponse{
			{
				Code:        "P0113",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}
		mockDeps.openAISvc.
			EXPECT().
			GetErrorCodesDescription(gomock.Eq("Toyota"), gomock.Eq("Camry"), gomock.Eq(req.ErrorCodes)).
			Return(chatGptResp, nil).
			AnyTimes()

		j, _ := json.Marshal(req)

		request := test.BuildRequest("POST", "/user/devices/"+ud.ID+"/error-codes", string(j))
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assert.Equal(t,
			"Too many error codes. Error codes list must be 100 or below in length.",
			string(body),
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ShouldErrorInvalidErrorCodes(t *testing.T) {

	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.QueryDeviceErrorCodes)

	t.Run("POST - get description for query codes", func(t *testing.T) {

		req := QueryDeviceErrorCodesReq{
			ErrorCodes: []string{"P0010:30", "P33333339"},
		}

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		mockDeps.deviceDefSvc.
			EXPECT().
			GetDeviceDefinitionByID(gomock.Any(), ud.DeviceDefinitionID).
			Return(&grpc.GetDeviceDefinitionItemResponse{
				Type: &grpc.DeviceType{
					Make:  "Toyota",
					Model: "Camry",
					Year:  2023,
				},
			}, nil).
			AnyTimes()

		chatGptResp := []services.ErrorCodesResponse{
			{
				Code:        "P0113",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}
		mockDeps.openAISvc.
			EXPECT().
			GetErrorCodesDescription(gomock.Eq("Toyota"), gomock.Eq("Camry"), gomock.Eq(req.ErrorCodes)).
			Return(chatGptResp, nil).
			AnyTimes()

		j, _ := json.Marshal(req)

		request := test.BuildRequest("POST", "/user/devices/"+ud.ID+"/error-codes", string(j))
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assert.Equal(t,
			"Invalid error code P33333339",
			string(body),
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ShouldErrorOnEmptyErrorCodes(t *testing.T) {

	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.QueryDeviceErrorCodes)

	t.Run("POST - get description for query codes", func(t *testing.T) {

		req := QueryDeviceErrorCodesReq{
			ErrorCodes: []string{},
		}

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		mockDeps.deviceDefSvc.
			EXPECT().
			GetDeviceDefinitionByID(gomock.Any(), ud.DeviceDefinitionID).
			Return(&grpc.GetDeviceDefinitionItemResponse{
				Type: &grpc.DeviceType{
					Make:  "Toyota",
					Model: "Camry",
					Year:  2023,
				},
			}, nil).
			AnyTimes()

		chatGptResp := []services.ErrorCodesResponse{
			{
				Code:        "P0113",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}
		mockDeps.openAISvc.
			EXPECT().
			GetErrorCodesDescription(gomock.Eq("Toyota"), gomock.Eq("Camry"), gomock.Eq(req.ErrorCodes)).
			Return(chatGptResp, nil).
			AnyTimes()

		j, _ := json.Marshal(req)

		request := test.BuildRequest("POST", "/user/devices/"+ud.ID+"/error-codes", string(j))
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assert.Equal(t,
			"No error codes provided",
			string(body),
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ShouldStoreErrorCodeResponse(t *testing.T) {

	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.QueryDeviceErrorCodes)

	t.Run("POST - get description for query codes", func(t *testing.T) {
		erCodeReq := []string{"P0017", "P0016"}
		req := QueryDeviceErrorCodesReq{
			ErrorCodes: erCodeReq,
		}

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		mockDeps.deviceDefSvc.
			EXPECT().
			GetDeviceDefinitionByID(gomock.Any(), ud.DeviceDefinitionID).
			Return(&grpc.GetDeviceDefinitionItemResponse{
				Type: &grpc.DeviceType{
					Make:  "Toyota",
					Model: "Camry",
					Year:  2023,
				},
			}, nil).
			AnyTimes()

		openAIResp := []services.ErrorCodesResponse{
			{
				Code:        "P0113",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}
		mockDeps.openAISvc.
			EXPECT().
			GetErrorCodesDescription(gomock.Eq("Toyota"), gomock.Eq("Camry"), gomock.Eq(req.ErrorCodes)).
			Return(openAIResp, nil).
			AnyTimes()

		j, _ := json.Marshal(req)

		request := test.BuildRequest("POST", "/user/devices/"+ud.ID+"/error-codes", string(j))
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		chatGptResp := QueryDeviceErrorCodesResponse{
			ErrorCodes: openAIResp,
		}
		chtJSON, err := json.Marshal(chatGptResp)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		assert.Equal(t,
			chtJSON,
			body,
		)

		errCodeResp, err := models.ErrorCodeQueries(
			models.ErrorCodeQueryWhere.UserDeviceID.EQ(ud.ID),
		).One(ctx, pdb.DBS().Reader)
		assert.NoError(t, err)

		ddd := null.JSONFrom([]byte(
			`[{"code": "P0113", "description": "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat."}]`,
		))

		assert.Equal(t, errCodeResp.CodesQueryResponse, ddd)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_GetUserDevicesErrorCodeQueries(t *testing.T) {
	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Get("/user/devices/:userDeviceID/error-codes", test.AuthInjectorTestHandler(testUserID), c.GetUserDeviceErrorCodeQueries)

	t.Run("GET - all saved error code response for current user devices", func(t *testing.T) {

		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		chatGptResp := []services.ErrorCodesResponse{
			{
				Code:        "P0017",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
			{
				Code:        "P0016",
				Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
			},
		}
		chtJSON, err := json.Marshal(chatGptResp)
		assert.NoError(t, err)

		currTime := time.Now().UTC().Truncate(time.Microsecond)
		erCodeQuery := models.ErrorCodeQuery{
			ID:                 ksuid.New().String(),
			UserDeviceID:       ud.ID,
			CodesQueryResponse: null.JSONFrom(chtJSON),
			CreatedAt:          currTime,
		}

		err = erCodeQuery.Insert(ctx, pdb.DBS().Writer, boil.Infer())
		assert.NoError(t, err)

		request := test.BuildRequest("GET", fmt.Sprintf("/user/devices/%s/error-codes", ud.ID), "")
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		resp := GetUserDeviceErrorCodeQueriesResponse{
			Queries: []GetUserDeviceErrorCodeQueriesResponseItem{
				{
					ErrorCodes:  chatGptResp,
					RequestedAt: currTime,
					ClearedAt:   erCodeQuery.ClearedAt.Ptr(),
				},
			},
		}

		expectedBody, err := json.Marshal(resp)
		assert.NoError(t, err)

		assert.JSONEq(t,
			string(expectedBody),
			string(body),
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ClearUserDeviceErrorCodeQuery(t *testing.T) {
	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes/clear", test.AuthInjectorTestHandler(testUserID), c.ClearUserDeviceErrorCodeQuery)

	t.Run("POST - clear last saved error code response for current user devices", func(t *testing.T) {
		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		testData := []struct {
			Codes      []string
			OpenAIResp []services.ErrorCodesResponse
		}{
			{
				Codes: []string{"P0017"},
				OpenAIResp: []services.ErrorCodesResponse{
					{
						Code:        "P0017",
						Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
					},
				},
			},
			{
				Codes: []string{"P0016"},
				OpenAIResp: []services.ErrorCodesResponse{
					{
						Code:        "P0016",
						Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
					},
				},
			},
		}

		for _, tData := range testData {
			chtJSON, err := json.Marshal(tData.OpenAIResp)
			assert.NoError(t, err)

			currTime := time.Now().UTC().Truncate(time.Microsecond)
			erCodeQuery := models.ErrorCodeQuery{
				ID:                 ksuid.New().String(),
				UserDeviceID:       ud.ID,
				CodesQueryResponse: null.JSONFrom(chtJSON),
				CreatedAt:          currTime,
			}

			err = erCodeQuery.Insert(ctx, pdb.DBS().Writer, boil.Infer())
			assert.NoError(t, err)
		}

		request := test.BuildRequest("POST", fmt.Sprintf("/user/devices/%s/error-codes/clear", ud.ID), "")
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		errCodeQuery, err := models.ErrorCodeQueries(
			models.ErrorCodeQueryWhere.ClearedAt.IsNotNull(),
			qm.OrderBy(models.ErrorCodeQueryColumns.CreatedAt+" DESC"),
			qm.Limit(1),
		).One(ctx, pdb.DBS().Reader)
		assert.NoError(t, err)

		currTime := errCodeQuery.ClearedAt.Time.UTC()

		assert.JSONEq(t,
			fmt.Sprintf(`{"errorCodes":%s, "clearedAt":"%s"}`, string(errCodeQuery.CodesQueryResponse.JSON), currTime.Format(time.RFC3339Nano)),
			string(body),
		)

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_ErrorOnAllErrorCodesCleared(t *testing.T) {
	mockDeps := createMockDependencies(t)
	defer mockDeps.mockCtrl.Finish()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	testUserID := "123123"
	c := NewUserDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &mockDeps.logger, mockDeps.deviceDefSvc, mockDeps.deviceDefIntSvc, &fakeEventService{}, mockDeps.scClient, mockDeps.scTaskSvc, mockDeps.teslaSvc, mockDeps.teslaTaskService, nil, nil, mockDeps.nhtsaService, mockDeps.autoPiIngest, mockDeps.deviceDefinitionIngest, mockDeps.autoPiTaskSvc, nil, nil, nil, nil, mockDeps.openAISvc, nil, nil, nil)
	app := fiber.New()
	app.Post("/user/devices/:userDeviceID/error-codes/clear", test.AuthInjectorTestHandler(testUserID), c.ClearUserDeviceErrorCodeQuery)

	t.Run("POST - clear last saved error code response for current user devices", func(t *testing.T) {
		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Toyota", "Camry", 2023, autoPiInteg)
		ud := test.SetupCreateUserDevice(t, testUserID, dd[0].DeviceDefinitionId, nil, "", pdb)

		testData := []struct {
			Codes      []string
			OpenAIResp []services.ErrorCodesResponse
		}{
			{
				Codes: []string{"P0017"},
				OpenAIResp: []services.ErrorCodesResponse{
					{
						Code:        "P0017",
						Description: "Engine Coolant Temperature Circuit Malfunction: This code indicates that the engine coolant temperature sensor is sending a signal that is outside of the expected range, which may cause the engine to run poorly or overheat.",
					},
				},
			},
		}

		for _, tData := range testData {
			chtJSON, err := json.Marshal(tData.OpenAIResp)
			assert.NoError(t, err)

			currTime := time.Now().UTC().Truncate(time.Microsecond)
			erCodeQuery := models.ErrorCodeQuery{
				ID:                 ksuid.New().String(),
				UserDeviceID:       ud.ID,
				CodesQueryResponse: null.JSONFrom(chtJSON),
				CreatedAt:          currTime,
				ClearedAt:          null.TimeFrom(currTime),
			}

			err = erCodeQuery.Insert(ctx, pdb.DBS().Writer, boil.Infer())
			assert.NoError(t, err)
		}

		request := test.BuildRequest("POST", fmt.Sprintf("/user/devices/%s/error-codes/clear", ud.ID), "")
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		assert.Equal(t, response.StatusCode, fiber.StatusBadRequest)
		assert.Equal(t, "all error codes already cleared", string(body))

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}
