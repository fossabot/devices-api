package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"strconv"

	"github.com/DIMO-Network/devices-api/internal/config"
	"github.com/DIMO-Network/devices-api/internal/database"
	"github.com/DIMO-Network/devices-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type NFTController struct {
	Settings *config.Settings
	DBS      func() *database.DBReaderWriter
	s3       *s3.Client
	log      *zerolog.Logger
}

// NewUserDevicesController constructor
func NewNFTController(
	settings *config.Settings,
	dbs func() *database.DBReaderWriter,
	logger *zerolog.Logger,
) NFTController {
	awscfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(settings.AWSRegion))
	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't create AWS config.")
	}
	s3Client := s3.NewFromConfig(awscfg)

	return NFTController{
		Settings: settings,
		DBS:      dbs,
		log:      logger,
		s3:       s3Client,
	}
}

// GetNFTMetadata godoc
// @Description  retrieves NFT metadata for a given tokenID
// @Tags         nfts
// @Produce      json
// @Success      200  {object}  controllers.NFTMetadataResp
// @Failure      404
// @Router       /nfts/:tokenID [get]
func (udc *NFTController) GetNFTMetadata(c *fiber.Ctx) error {
	tis := c.Params("tokenID")
	ti, ok := new(big.Int).SetString(tis, 10)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tis))
	}

	tid := types.NewNullDecimal(new(decimal.Big).SetBigMantScale(ti, 0))

	mr, err := models.MintRequests(
		models.MintRequestWhere.TokenID.EQ(tid),
		qm.Load(qm.Rels(models.MintRequestRels.UserDevice, models.UserDeviceRels.DeviceDefinition, models.DeviceDefinitionRels.DeviceMake)),
	).One(c.Context(), udc.DBS().Writer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "NFT not found.")
		}
		return opaqueInternalError
	}

	description := fmt.Sprintf("%s %s %d", mr.R.UserDevice.R.DeviceDefinition.R.DeviceMake.Name, mr.R.UserDevice.R.DeviceDefinition.Model, mr.R.UserDevice.R.DeviceDefinition.Year)

	var name string
	if mr.R.UserDevice.Name.Valid {
		name = mr.R.UserDevice.Name.String
	} else {
		name = description
	}

	return c.JSON(NFTMetadataResp{
		Name:        name,
		Description: description,
		Image:       fmt.Sprintf("%s/v1/nfts/%s/image", udc.Settings.DeploymentBaseURL, ti),
		Attributes: []NFTAttribute{
			{TraitType: "Make", Value: mr.R.UserDevice.R.DeviceDefinition.R.DeviceMake.Name},
			{TraitType: "Model", Value: mr.R.UserDevice.R.DeviceDefinition.Model},
			{TraitType: "Year", Value: strconv.Itoa(int(mr.R.UserDevice.R.DeviceDefinition.Year))},
		},
	})
}

type NFTMetadataResp struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Attributes  []NFTAttribute `json:"attributes"`
}

type NFTAttribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

// GetNFTImage godoc
// @Description  retrieves NFT metadata for a given tokenID
// @Tags         nfts
// @Produce      png
// @Router       /nfts/:tokenID/image [get]
func (udc *NFTController) GetNFTImage(c *fiber.Ctx) error {
	tis := c.Params("tokenID")
	ti, ok := new(big.Int).SetString(tis, 10)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tis))
	}

	tid := types.NewNullDecimal(new(decimal.Big).SetBigMantScale(ti, 0))

	mr, err := models.MintRequests(
		models.MintRequestWhere.TokenID.EQ(tid),
	).One(c.Context(), udc.DBS().Writer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "NFT not found.")
		}
		return opaqueInternalError
	}

	s3o, err := udc.s3.GetObject(c.Context(), &s3.GetObjectInput{
		Bucket: aws.String(udc.Settings.NFTS3Bucket),
		Key:    aws.String(mr.ID + ".png"),
	})
	if err != nil {
		udc.log.Err(err).Msg("Failure communicating with S3.")
		return opaqueInternalError
	}

	c.Set("Content-Type", "image/png")
	return c.SendStream(s3o.Body)
}