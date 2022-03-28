package controllers

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DIMO-Network/devices-api/internal/config"
	"github.com/DIMO-Network/devices-api/internal/services"
	mock_services "github.com/DIMO-Network/devices-api/internal/services/mocks"
	"github.com/DIMO-Network/devices-api/internal/test"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
)

// integration tests using embedded pgsql, must be run in order
func TestDevicesController(t *testing.T) {
	// arrange global db and route setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

	ctx := context.Background()
	pdb := test.GetDBConnection(ctx)

	nhtsaSvc := mock_services.NewMockINHTSAService(mockCtrl)
	deviceDefSvc := mock_services.NewMockIDeviceDefinitionService(mockCtrl)
	c := NewDevicesController(&config.Settings{Port: "3000"}, pdb.DBS, &logger, nhtsaSvc, deviceDefSvc)
	// routes
	app := fiber.New()
	app.Get("/device-definitions/all", c.GetAllDeviceMakeModelYears)
	app.Get("/device-definitions/:id", c.GetDeviceDefinitionByID)
	app.Get("/device-definitions/:id/integrations", c.GetDeviceIntegrationsByID)

	dbMake := test.SetupCreateMake(t, "TESLA", pdb)
	dbDeviceDef := test.SetupCreateDeviceDefinition(t, dbMake, "MODEL Y", 2020, pdb)
	fmt.Println("created device def id: " + dbDeviceDef.ID)
	createdID := dbDeviceDef.ID

	t.Run("GET - device definition by id, including autopi integration", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/device-definitions/"+createdID, nil)
		response, _ := app.Test(request)
		body, _ := ioutil.ReadAll(response.Body)
		// assert
		assert.Equal(t, 200, response.StatusCode)

		v := gjson.GetBytes(body, "deviceDefinition")
		var dd services.DeviceDefinition
		err := json.Unmarshal([]byte(v.Raw), &dd)
		assert.NoError(t, err)
		assert.Equal(t, createdID, dd.DeviceDefinitionID)
		if assert.True(t, len(dd.CompatibleIntegrations) >= 2, "should be atleast 2 integrations for autopi") {
			assert.Equal(t, services.AutoPiVendor, dd.CompatibleIntegrations[0].Vendor)
			assert.Equal(t, "Americas", dd.CompatibleIntegrations[0].Region)
			assert.Equal(t, services.AutoPiVendor, dd.CompatibleIntegrations[1].Vendor)
			assert.Equal(t, "Europe", dd.CompatibleIntegrations[1].Region)
		}
	})
	t.Run("GET - device definition by id does not add autopi compatibility for old vehicle", func(t *testing.T) {
		dbDdOldCar := test.SetupCreateDeviceDefinition(t, dbMake, "Oldie", 1999, pdb)
		request, _ := http.NewRequest("GET", "/device-definitions/"+dbDdOldCar.ID, nil)
		response, _ := app.Test(request)
		body, _ := ioutil.ReadAll(response.Body)
		// assert
		assert.Equal(t, 200, response.StatusCode)
		v := gjson.GetBytes(body, "deviceDefinition")
		var dd services.DeviceDefinition
		err := json.Unmarshal([]byte(v.Raw), &dd)
		assert.NoError(t, err)
		assert.Equal(t, dbDdOldCar.ID, dd.DeviceDefinitionID)
		assert.Len(t, dd.CompatibleIntegrations, 0, "vehicles before 2020 should not auto inject autopi integrations")
	})
	t.Run("GET - device integrations by id", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/device-definitions/"+createdID+"/integrations", nil)
		response, _ := app.Test(request)
		body, _ := ioutil.ReadAll(response.Body)
		// assert
		assert.Equal(t, 200, response.StatusCode)
		v := gjson.GetBytes(body, "compatibleIntegrations")
		var dc []services.DeviceCompatibility
		err := json.Unmarshal([]byte(v.Raw), &dc)
		assert.NoError(t, err)
		if assert.True(t, len(dc) >= 2, "should be atleast 2 integrations for autopi") {
			assert.Equal(t, services.AutoPiVendor, dc[0].Vendor)
			assert.Equal(t, "Americas", dc[0].Region)
			assert.Equal(t, services.AutoPiVendor, dc[1].Vendor)
			assert.Equal(t, "Europe", dc[1].Region)
		}
	})
	t.Run("GET 400 - device definition by id invalid", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/device-definitions/caca", nil)
		response, _ := app.Test(request)
		// assert
		assert.Equal(t, 400, response.StatusCode)
	})
	t.Run("GET 400 - device definition integrations invalid", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/device-definitions/caca/integrations", nil)
		response, _ := app.Test(request)
		// assert
		assert.Equal(t, 400, response.StatusCode)
	})
	t.Run("GET - all make model years as a tree", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/device-definitions/all", nil)
		response, _ := app.Test(request)
		body, _ := ioutil.ReadAll(response.Body)
		// assert
		assert.Equal(t, 200, response.StatusCode)
		v := gjson.GetBytes(body, "makes")
		var mmy []DeviceMMYRoot
		err := json.Unmarshal([]byte(v.Raw), &mmy)
		assert.NoError(t, err)
		if assert.True(t, len(mmy) >= 1, "should be at least one device definition") {
			assert.Equal(t, "TESLA", mmy[0].Make)
			assert.Equal(t, "MODEL Y", mmy[0].Models[0].Model)
			assert.Equal(t, int16(2020), mmy[0].Models[0].Years[0].Year)
			assert.Equal(t, createdID, mmy[0].Models[0].Years[0].DeviceDefinitionID)
		}
	})
}

func TestNewDeviceDefinitionFromDatabase(t *testing.T) {
	dbMake := &models.DeviceMake{
		ID:   ksuid.New().String(),
		Name: "Mercedes",
	}
	dbDevice := models.DeviceDefinition{
		ID:           "123",
		DeviceMakeID: dbMake.ID,
		Model:        "R500",
		Year:         2020,
		Metadata:     null.JSONFrom([]byte(`{"vehicle_info": {"fuel_type": "gas", "driven_wheels": "4", "number_of_doors":"5" } }`)),
	}
	ds := models.DeviceStyle{
		SubModel:           "AMG",
		Name:               "C63 AMG",
		DeviceDefinitionID: dbDevice.ID,
	}
	di := models.DeviceIntegration{
		DeviceDefinitionID: "123",
		IntegrationID:      "123",
		CreatedAt:          time.Time{},
		UpdatedAt:          time.Time{},
	}
	di.R = di.R.NewStruct()
	di.R.Integration = &models.Integration{
		ID:     "123",
		Type:   "Hardware",
		Style:  "Addon",
		Vendor: "Autopi",
	}
	dbDevice.R = dbDevice.R.NewStruct()
	dbDevice.R.DeviceMake = dbMake
	dbDevice.R.DeviceIntegrations = append(dbDevice.R.DeviceIntegrations, &di)
	dbDevice.R.DeviceStyles = append(dbDevice.R.DeviceStyles, &ds)
	dd, err := NewDeviceDefinitionFromDatabase(&dbDevice)

	assert.NoError(t, err)
	assert.Equal(t, "123", dd.DeviceDefinitionID)
	assert.Equal(t, "gas", dd.VehicleInfo.FuelType)
	assert.Equal(t, "4", dd.VehicleInfo.DrivenWheels)
	assert.Equal(t, "5", dd.VehicleInfo.NumberOfDoors)
	assert.Equal(t, "Vehicle", dd.Type.Type)
	assert.Equal(t, 2020, dd.Type.Year)
	assert.Equal(t, "Mercedes", dd.Type.Make)
	assert.Equal(t, "R500", dd.Type.Model)
	assert.Contains(t, dd.Type.SubModels, "AMG")

	assert.Len(t, dd.CompatibleIntegrations, 1)
	assert.Equal(t, "Autopi", dd.CompatibleIntegrations[0].Vendor)
}

func TestNewDeviceDefinitionFromDatabase_Error(t *testing.T) {
	dbDevice := models.DeviceDefinition{
		ID:       "123",
		Model:    "R500",
		Year:     2020,
		Metadata: null.JSONFrom([]byte(`{"vehicle_info": {"fuel_type": "gas", "driven_wheels": "4", "number_of_doors":"5" } }`)),
	}
	dbDevice.R = dbDevice.R.NewStruct()
	_, err := NewDeviceDefinitionFromDatabase(&dbDevice)
	assert.Error(t, err)

	dbDevice.R = nil
	_, err = NewDeviceDefinitionFromDatabase(&dbDevice)
	assert.Error(t, err)
}
