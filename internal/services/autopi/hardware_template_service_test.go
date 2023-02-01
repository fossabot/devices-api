package autopi

import (
	"encoding/json"
	"github.com/DIMO-Network/devices-api/internal/services"
	"strconv"
	"testing"

	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/devices-api/internal/constants"
	"github.com/DIMO-Network/devices-api/internal/test"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
)

type HardwareTemplateServiceTestSuite struct {
	suite.Suite
	hardwareTemplateService HardwareTemplateService
}

func (s *HardwareTemplateServiceTestSuite) SetupSuite() {
	mockCtrl := gomock.NewController(s.T())
	defer mockCtrl.Finish()
	s.hardwareTemplateService = NewHardwareTemplateService()
}

func (s *HardwareTemplateServiceTestSuite) TearDownTest() {
}

func (s *HardwareTemplateServiceTestSuite) TearDownSuite() {
}

func TestHardwareTemplateServiceTestSuite(t *testing.T) {
	suite.Run(t, new(HardwareTemplateServiceTestSuite))
}

func (s *HardwareTemplateServiceTestSuite) Test_HardwareTemplateService() {
	type tableTestCases struct {
		description string
		expected    string
		ud          *models.UserDevice
		dd          *ddgrpc.GetDeviceDefinitionItemResponse
		integ       *ddgrpc.Integration
	}
	const (
		tIDIntegrationDefault = "10"
		tIDDeviceStyle        = "11"
		tIDDeviceStyleFromUD  = "111"
		tIDDeviceDef          = "12"
		tIDDeviceMake         = "13"
		tIDBEVPowertrainUD    = "14"
	)
	def, _ := strconv.Atoi(tIDIntegrationDefault)
	bev, _ := strconv.Atoi(tIDBEVPowertrainUD)
	integration := test.BuildIntegrationDefaultGRPC(constants.AutoPiVendor, def, bev, true)
	integrationWithoutAutoPiPowertrainTemplate := test.BuildIntegrationDefaultGRPC(constants.AutoPiVendor, def, 0, false)

	ddWithTID := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "F150", 2020, integration)[0]
	ddWithDeviceStyleTID := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "F150", 2020, integration)[0]
	ddWithMakeTID := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "F150", 2020, integration)[0]
	ddNoTIDs := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "F150", 2020, integration)[0]
	ddWithDeviceStyleInUD := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "F150", 2020, integration)[0]

	ddWithTID.HardwareTemplateId = tIDDeviceDef
	ddWithDeviceStyleTID.DeviceStyles = append(ddWithDeviceStyleTID.DeviceStyles, &ddgrpc.DeviceStyle{
		Id:                 ksuid.New().String(),
		HardwareTemplateId: tIDDeviceStyle,
	})
	ddWithDeviceStyleInUD.DeviceStyles = append(ddWithDeviceStyleInUD.DeviceStyles, &ddgrpc.DeviceStyle{
		Id:                 ksuid.New().String(),
		HardwareTemplateId: tIDDeviceStyleFromUD,
	})
	ddWithMakeTID.Make.HardwareTemplateId = tIDDeviceMake

	pt := services.BEV
	udmdBEVPT := services.UserDeviceMetadata{
		PowertrainType: &pt,
	}
	udmdBEVPTjson, _ := json.Marshal(udmdBEVPT)

	for _, scenario := range []tableTestCases{
		{
			description: "Should get hardware template id from style id in User Device",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddWithDeviceStyleInUD.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
				DeviceStyleID:      null.StringFrom(ddWithDeviceStyleInUD.DeviceStyles[0].Id),
			},
			integ:    integration,
			dd:       ddWithDeviceStyleInUD,
			expected: tIDDeviceStyleFromUD,
		},
		{
			description: "Should get hardware template id from DD",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddWithTID.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
			},
			integ:    integration,
			dd:       ddWithTID,
			expected: tIDDeviceDef,
		},
		{
			description: "Should get template id from DD with styles but no style id in UD",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddWithDeviceStyleTID.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
			},
			integ:    integration,
			dd:       ddWithDeviceStyleTID,
			expected: tIDDeviceStyle,
		},
		{
			description: "Should get hardware template id from Device Make in DD",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddWithMakeTID.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
			},
			integ:    integration,
			dd:       ddWithMakeTID,
			expected: tIDDeviceMake,
		},
		{
			description: "Should get hardware template id from AutoPi integration AutoPiPowertrainTemplate in UD",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddNoTIDs.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
				Metadata:           null.JSONFrom(udmdBEVPTjson),
			},
			integ:    integration,
			dd:       ddNoTIDs,
			expected: tIDBEVPowertrainUD,
		},
		{
			description: "Should get hardware template id from AutoPi DefaultTemplate",
			ud: &models.UserDevice{
				ID:                 ksuid.New().String(),
				UserID:             "testUserID",
				DeviceDefinitionID: ddNoTIDs.DeviceDefinitionId,
				CountryCode:        null.StringFrom("USA"),
				Name:               null.StringFrom("Chungus"),
			},
			integ:    integrationWithoutAutoPiPowertrainTemplate,
			dd:       ddNoTIDs,
			expected: tIDIntegrationDefault,
		},
	} {
		s.T().Run(scenario.description, func(t *testing.T) {
			id, _ := s.hardwareTemplateService.GetTemplateID(scenario.ud, scenario.dd, scenario.integ)
			assert.Equal(t, scenario.expected, id)
		})
	}
}
