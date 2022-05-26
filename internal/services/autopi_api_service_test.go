package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/DIMO-Network/devices-api/internal/database"
	"github.com/DIMO-Network/devices-api/internal/test"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AutoPiAPIServiceTestSuite struct {
	suite.Suite
	pdb       database.DbStore
	container testcontainers.Container
	ctx       context.Context
}

// SetupSuite starts container db
func (s *AutoPiAPIServiceTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.pdb, s.container = test.StartContainerDatabase(s.ctx, s.T(), migrationsDirRelPath)
}

//TearDownTest after each test truncate tables
func (s *AutoPiAPIServiceTestSuite) TearDownTest() {
	test.TruncateTables(s.pdb.DBS().Writer.DB, s.T())
}

//TearDownSuite cleanup at end by terminating container
func (s *AutoPiAPIServiceTestSuite) TearDownSuite() {
	fmt.Printf("shutting down postgres at with session: %s \n", s.container.SessionID())
	if err := s.container.Terminate(s.ctx); err != nil {
		s.T().Fatal(err)
	}
}

func TestAutoPiApiServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AutoPiAPIServiceTestSuite))
}

func (s *AutoPiAPIServiceTestSuite) TestFindUserDeviceAutoPiIntegration() {
	// arrange some data
	const testUserID = "123123"
	const autoPiDeviceID = "321"
	autoPiUnitID := "456"
	apInt := test.SetupCreateAutoPiIntegration(s.T(), 10, nil, s.pdb)
	scInt := test.SetupCreateSmartCarIntegration(s.T(), s.pdb)
	dm := test.SetupCreateMake(s.T(), "Tesla", s.pdb)
	dd := test.SetupCreateDeviceDefinition(s.T(), dm, "Model 3", 2020, s.pdb)
	test.SetupCreateDeviceIntegration(s.T(), dd, apInt, s.pdb)
	test.SetupCreateDeviceIntegration(s.T(), dd, scInt, s.pdb)
	ud := test.SetupCreateUserDevice(s.T(), testUserID, dd, nil, s.pdb)
	// now create the api ints
	scUdai := &models.UserDeviceAPIIntegration{
		UserDeviceID:  ud.ID,
		IntegrationID: scInt.ID,
		Status:        models.UserDeviceAPIIntegrationStatusActive,
		ExternalID:    null.StringFrom("423324"),
	}
	err := scUdai.Insert(s.ctx, s.pdb.DBS().Writer, boil.Infer())
	assert.NoError(s.T(), err)
	amd := UserDeviceAPIIntegrationsMetadata{
		AutoPiUnitID: &autoPiUnitID,
	}
	apUdai := &models.UserDeviceAPIIntegration{
		UserDeviceID:  ud.ID,
		IntegrationID: apInt.ID,
		Status:        models.UserDeviceAPIIntegrationStatusActive,
		ExternalID:    null.StringFrom(autoPiDeviceID),
	}
	_ = apUdai.Metadata.Marshal(amd)
	err = apUdai.Insert(s.ctx, s.pdb.DBS().Writer, boil.Infer())
	assert.NoError(s.T(), err)
	// act  now call the method
	udIntegration, metadata, err := FindUserDeviceAutoPiIntegration(s.ctx, s.pdb.DBS().Writer, ud.ID, testUserID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), udIntegration, "expected user_device_api_integration not be nil")
	assert.NotNilf(s.T(), metadata, "expected metadata not be nil")
	assert.Equal(s.T(), ud.ID, udIntegration.UserDeviceID)
	assert.Equal(s.T(), apInt.ID, udIntegration.IntegrationID)
	assert.Equal(s.T(), autoPiDeviceID, udIntegration.ExternalID.String)
	// check some values
	test.TruncateTables(s.pdb.DBS().Writer.DB, s.T())
}

func (s *AutoPiAPIServiceTestSuite) TestAppendAutoPiCompatibility() {
	dm := test.SetupCreateMake(s.T(), "Ford", s.pdb)
	dd := test.SetupCreateDeviceDefinition(s.T(), dm, "MachE", 2020, s.pdb)
	var dcs []DeviceCompatibility
	compatibility, err := AppendAutoPiCompatibility(s.ctx, dcs, dd.ID, s.pdb.DBS().Writer)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), compatibility, 2)
	all, err := models.DeviceIntegrations().All(s.ctx, s.pdb.DBS().Reader)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), all, 2)

	test.TruncateTables(s.pdb.DBS().Writer.DB, s.T())
}
