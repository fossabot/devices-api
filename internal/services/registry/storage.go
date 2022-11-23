package registry

import (
	"context"

	"github.com/DIMO-Network/devices-api/internal/database"
	"github.com/DIMO-Network/devices-api/internal/services/autopi"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type Storage interface {
	HandleUpdate(ctx context.Context, data *ceData) error
}

type S struct {
	ABI    *abi.ABI
	DB     func() *database.DBReaderWriter
	Logger *zerolog.Logger
	ap     *autopi.Integration
}

func (s *S) HandleUpdate(ctx context.Context, data *ceData) error {
	logger := s.Logger.With().
		Str("requestId", data.RequestID).
		Str("status", data.Type).
		Str("hash", data.Transaction.Hash).
		Logger()

	logger.Info().Msg("Got transaction status.")

	mtr, err := models.MetaTransactionRequests(
		models.MetaTransactionRequestWhere.ID.EQ(data.RequestID),
		// This is really ugly. We should probably link back to the type instead of doing this.
		qm.Load(models.MetaTransactionRequestRels.MintRequestVehicleNFT),
		qm.Load(models.MetaTransactionRequestRels.ClaimMetaTransactionRequestAutopiUnit),
		qm.Load(models.MetaTransactionRequestRels.PairRequestAutopiUnit),
		qm.Load(models.MetaTransactionRequestRels.UnpairRequestAutopiUnit),
	).One(context.Background(), s.DB().Reader)
	if err != nil {
		return err
	}

	mtr.Status = data.Type
	mtr.Hash = null.BytesFrom(common.FromHex(data.Transaction.Hash))

	_, err = mtr.Update(ctx, s.DB().Writer, boil.Infer())
	if err != nil {
		return err
	}

	if mtr.Status != models.MetaTransactionRequestStatusConfirmed {
		return nil
	}

	vehicleMintedEvent := s.ABI.Events["VehicleNodeMinted"]
	deviceClaimedEvent := s.ABI.Events["AftermarketDeviceClaimed"]
	devicePairedEvent := s.ABI.Events["AftermarketDevicePaired"]
	deviceUnpairedEvent := s.ABI.Events["AftermarketDeviceUnpaired"]

	switch {
	case mtr.R.MintRequestVehicleNFT != nil:
		for _, l1 := range data.Transaction.Logs {
			l2 := convertLog(&l1)
			if l2.Topics[0] == vehicleMintedEvent.ID {
				out := new(RegistryVehicleNodeMinted)
				err := s.parseLog(out, vehicleMintedEvent, *l2)
				if err != nil {
					return err
				}

				mtr.R.MintRequestVehicleNFT.TokenID = types.NewNullDecimal(new(decimal.Big).SetBigMantScale(out.TokenId, 0))
				mtr.R.MintRequestVehicleNFT.OwnerAddress = null.BytesFrom(out.Owner.Bytes())
				_, err = mtr.R.MintRequestVehicleNFT.Update(ctx, s.DB().Writer, boil.Infer())
				if err != nil {
					return err
				}

				logger.Info().Str("userDeviceId", mtr.R.MintRequestVehicleNFT.UserDeviceID.String).Msg("Vehicle minted.")
			}
		}
		// Other soon.

	case mtr.R.ClaimMetaTransactionRequestAutopiUnit != nil:
		for _, l1 := range data.Transaction.Logs {
			l2 := convertLog(&l1)
			if l2.Topics[0] == deviceClaimedEvent.ID {
				out := new(RegistryAftermarketDeviceClaimed)
				err := s.parseLog(out, deviceClaimedEvent, *l2)
				if err != nil {
					return err
				}

				mtr.R.ClaimMetaTransactionRequestAutopiUnit.OwnerAddress = null.BytesFrom(out.Owner[:])
				_, err = mtr.R.ClaimMetaTransactionRequestAutopiUnit.Update(ctx, s.DB().Writer, boil.Infer())
				if err != nil {
					return err
				}

				logger.Info().Str("autoPiTokenId", mtr.R.ClaimMetaTransactionRequestAutopiUnit.TokenID.String()).Str("owner", out.Owner.String()).Msg("Device claimed.")
			}
		}
	case mtr.R.PairRequestAutopiUnit != nil:
		for _, l1 := range data.Transaction.Logs {
			l2 := convertLog(&l1)
			if l2.Topics[0] == devicePairedEvent.ID {
				out := new(RegistryAftermarketDevicePaired)
				err := s.parseLog(out, devicePairedEvent, *l2)
				if err != nil {
					return err
				}

				return s.ap.Pair(ctx, out.AftermarketDeviceNode, out.VehicleNode)
			}
		}
	case mtr.R.UnpairRequestAutopiUnit != nil:
		for _, l1 := range data.Transaction.Logs {
			l2 := convertLog(&l1)
			if l2.Topics[0] == deviceUnpairedEvent.ID {
				out := new(RegistryAftermarketDeviceUnpaired)
				err := s.parseLog(out, deviceUnpairedEvent, *l2)
				if err != nil {
					return err
				}

				mtr.R.UnpairRequestAutopiUnit.PairRequestID = null.String{}
				_, err = mtr.R.UnpairRequestAutopiUnit.Update(ctx, s.DB().Writer, boil.Infer())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *S) parseLog(out any, event abi.Event, log eth_types.Log) error {
	if len(log.Data) > 0 {
		err := s.ABI.UnpackIntoInterface(out, event.Name, log.Data)
		if err != nil {
			return err
		}
	}

	var indexed abi.Arguments
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}

	err := abi.ParseTopics(out, indexed, log.Topics[1:])
	if err != nil {
		return err
	}

	return nil
}

func convertLog(logIn *ceLog) *eth_types.Log {
	topics := make([]common.Hash, len(logIn.Topics))
	for i, t := range logIn.Topics {
		topics[i] = common.HexToHash(t)
	}

	data := common.FromHex(logIn.Data)

	return &eth_types.Log{
		Topics: topics,
		Data:   data,
	}
}

func NewStorage(db func() *database.DBReaderWriter, logger *zerolog.Logger, ap *autopi.Integration) (Storage, error) {
	abi, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &S{
		ABI:    abi,
		DB:     db,
		Logger: logger,
		ap:     ap,
	}, nil
}
