package fingerprint

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/DIMO-Network/devices-api/internal/services/issuer"
	"github.com/DIMO-Network/devices-api/internal/test"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type ConsumerTestSuite struct {
	suite.Suite
	pdb       db.Store
	container testcontainers.Container
	ctx       context.Context
	mockCtrl  *gomock.Controller
	topic     string
	iss       *issuer.Issuer
	cons      *Consumer
}

const migrationsDirRelPath = "../../../migrations"

// SetupSuite starts container db
func (s *ConsumerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.pdb, s.container = test.StartContainerDatabase(context.Background(), s.T(), migrationsDirRelPath)
	s.mockCtrl = gomock.NewController(s.T())
	s.topic = "topic.fingerprint"

	pk, err := base64.RawURLEncoding.DecodeString("2pN28-5VmEavX46XWszjasN0kx4ha3wQ6w6hGqD8o0k")
	require.NoError(s.T(), err)

	gitSha1 := os.Getenv("GIT_SHA1")
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Str("git-sha1", gitSha1).
		Logger()

	iss, err := issuer.New(issuer.Config{
		PrivateKey:        pk,
		ChainID:           big.NewInt(137),
		VehicleNFTAddress: common.HexToAddress("00f1"),
		DBS:               s.pdb,
	},
		&logger)
	s.iss = iss
	s.cons = &Consumer{
		logger: &logger,
		iss:    iss,
		DBS:    s.pdb,
	}

	s.Require().NoError(err)
}

// TearDownSuite cleanup at end by terminating container
func (s *ConsumerTestSuite) TearDownSuite() {
	fmt.Printf("shutting down postgres at with session: %s \n", s.container.SessionID())
	if err := s.container.Terminate(s.ctx); err != nil {
		s.T().Fatal(err)
	}
	s.mockCtrl.Finish()
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (s *ConsumerTestSuite) TestVinCredentialerHandler() {
	deviceID := ksuid.New().String()
	ownerAddress := null.BytesFrom(common.Hex2Bytes("ab8438a18d83d41847dffbdc6101d37c69c9a2fc"))
	vin := "1G6AL1RY2K0111939"
	ctx := context.Background()
	tokenID := big.NewInt(3)
	userDeviceID := "userDeviceID1"
	mtxReq := ksuid.New().String()
	deiceDefID := "deviceDefID"
	claimID := "claimID1"
	signature := "0xa4438e5cb667dc63ebd694167ae3ad83585f2834c9b04895dd890f805c4c459a024ed9df1b03872536b4ac0c7720d02cb787884a093cfcde5c3bd7f94657e30c1b"

	// tables used in tests
	aftermarketDevice := models.AftermarketDevice{
		UserID:          null.StringFrom("SomeID"),
		OwnerAddress:    ownerAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		TokenID:         types.NewNullDecimal(new(decimal.Big).SetBigMantScale(big.NewInt(13), 0)),
		VehicleTokenID:  types.NewNullDecimal(new(decimal.Big).SetBigMantScale(tokenID, 0)),
		Beneficiary:     null.BytesFrom(common.BytesToAddress([]byte{uint8(1)}).Bytes()),
		EthereumAddress: ownerAddress,
	}

	userDevice := models.UserDevice{
		ID:                 deviceID,
		UserID:             userDeviceID,
		DeviceDefinitionID: deiceDefID,
		VinConfirmed:       true,
		VinIdentifier:      null.StringFrom(vin),
	}

	metaTx := models.MetaTransactionRequest{
		ID:     mtxReq,
		Status: models.MetaTransactionRequestStatusConfirmed,
	}

	credential := models.VerifiableCredential{
		ClaimID:        claimID,
		Credential:     []byte{},
		ExpirationDate: time.Now().AddDate(0, 0, 7),
	}

	nft := models.VehicleNFT{
		MintRequestID: mtxReq,
		UserDeviceID:  null.StringFrom(deviceID),
		Vin:           vin,
		TokenID:       types.NewNullDecimal(new(decimal.Big).SetBigMantScale(tokenID, 0)),
		OwnerAddress:  ownerAddress,
		ClaimID:       null.StringFrom(claimID),
	}

	rawMsg, err := json.Marshal(struct {
		Vin string
	}{
		Vin: vin,
	})
	require.NoError(s.T(), err)

	cases := []struct {
		Name              string
		ReturnsError      bool
		ExpectedResponse  string
		UserDeviceTable   models.UserDevice
		MetaTxTable       models.MetaTransactionRequest
		VCTable           models.VerifiableCredential
		VehicleNFT        models.VehicleNFT
		AftermarketDevice models.AftermarketDevice
	}{
		{
			Name:             "No corresponding aftermarket device for address",
			ReturnsError:     true,
			ExpectedResponse: "sql: no rows in result set",
		},
		{
			Name:              "active credential",
			ReturnsError:      false,
			UserDeviceTable:   userDevice,
			MetaTxTable:       metaTx,
			VCTable:           credential,
			VehicleNFT:        nft,
			AftermarketDevice: aftermarketDevice,
		},
		{
			Name:            "inactive credential",
			ReturnsError:    false,
			UserDeviceTable: userDevice,
			MetaTxTable:     metaTx,
			VCTable: models.VerifiableCredential{
				ClaimID:        claimID,
				Credential:     []byte{},
				ExpirationDate: time.Now().AddDate(0, 0, -10),
			},
			VehicleNFT:        nft,
			AftermarketDevice: aftermarketDevice,
		},
		{
			Name:            "invalid token id",
			ReturnsError:    false,
			UserDeviceTable: userDevice,
			MetaTxTable:     metaTx,
			VCTable:         credential,
			VehicleNFT:      nft,
			AftermarketDevice: models.AftermarketDevice{
				UserID:          null.StringFrom("SomeID"),
				OwnerAddress:    ownerAddress,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				TokenID:         types.NewNullDecimal(new(decimal.Big).SetBigMantScale(big.NewInt(13), 0)),
				Beneficiary:     null.BytesFrom(common.BytesToAddress([]byte{uint8(1)}).Bytes()),
				EthereumAddress: ownerAddress,
			},
		},
	}

	for _, c := range cases {
		s.T().Run(c.Name, func(t *testing.T) {
			err := c.UserDeviceTable.Insert(ctx, s.pdb.DBS().Writer, boil.Infer())
			require.NoError(t, err)

			err = c.MetaTxTable.Insert(ctx, s.pdb.DBS().Writer, boil.Infer())
			require.NoError(t, err)

			err = c.VCTable.Insert(ctx, s.pdb.DBS().Reader, boil.Infer())
			require.NoError(t, err)

			err = c.VehicleNFT.Insert(ctx, s.pdb.DBS().Writer, boil.Infer())
			require.NoError(t, err)

			err = c.AftermarketDevice.Insert(ctx, s.pdb.DBS().Writer, boil.Infer())
			require.NoError(t, err)

			err = s.cons.Handle(s.ctx, &Event{
				CloudEvent: shared.CloudEvent[json.RawMessage]{
					Data:    rawMsg,
					Time:    time.Now(),
					ID:      deviceID,
					Subject: common.Bytes2Hex(ownerAddress.Bytes),
				},
				Signature: signature,
			})

			if c.ReturnsError {
				assert.ErrorContains(t, err, c.ExpectedResponse)
			} else {
				require.NoError(t, err)
			}

			test.TruncateTables(s.pdb.DBS().Writer.DB, t)
		})
	}

}