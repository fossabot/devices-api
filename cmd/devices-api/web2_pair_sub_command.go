package main

import (
	"context"
	"flag"
	"os"

	"math/big"

	"github.com/DIMO-Network/devices-api/internal/config"
	"github.com/DIMO-Network/devices-api/internal/services"
	"github.com/DIMO-Network/devices-api/internal/services/autopi"
	"github.com/DIMO-Network/shared/db"
	"github.com/google/subcommands"
	"github.com/rs/zerolog"
)

type web2PairCmd struct {
	logger    zerolog.Logger
	settings  config.Settings
	pdb       db.Store
	container dependencyContainer
}

func (*web2PairCmd) Name() string     { return "web2-pair" }
func (*web2PairCmd) Synopsis() string { return "web2-pair args to stdout." }
func (*web2PairCmd) Usage() string {
	return `web2-pair [] <some text>:
	web2-pair args.
  `
}

// nolint
func (p *web2PairCmd) SetFlags(f *flag.FlagSet) {

}

func (p *web2PairCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	producer := p.container.getKafkaProducer()
	ddSvc := p.container.getDeviceDefinitionService()

	if len(os.Args[2:]) != 2 {
		p.logger.Fatal().Msg("Requires aftermarket_token_id vehicle_token_id")
	}

	amToken, ok := new(big.Int).SetString(os.Args[2], 10)
	if !ok {
		p.logger.Fatal().Msgf("Couldn't parse aftermarket_token_id %q", os.Args[2])
	}

	vToken, ok := new(big.Int).SetString(os.Args[3], 10)
	if !ok {
		p.logger.Fatal().Msgf("Couldn't parse vehicle_token_id %q", os.Args[3])
	}

	p.logger.Info().Msgf("Attempting to web2 pair am device %s to vehicle %s.", amToken, vToken)

	autoPiSvc := services.NewAutoPiAPIService(&p.settings, p.pdb.DBS)
	autoPiTaskService := services.NewAutoPiTaskService(&p.settings, autoPiSvc, p.pdb.DBS, p.logger)
	autoPiIngest := services.NewIngestRegistrar(producer)
	eventService := services.NewEventService(&p.logger, &p.settings, producer)
	deviceDefinitionRegistrar := services.NewDeviceDefinitionRegistrar(producer, &p.settings)
	hardwareTemplateService := autopi.NewHardwareTemplateService(autoPiSvc, p.pdb.DBS, &p.logger)

	i := autopi.NewIntegration(p.pdb.DBS, ddSvc, autoPiSvc, autoPiTaskService, autoPiIngest, eventService, deviceDefinitionRegistrar, hardwareTemplateService, &p.logger)

	err := i.Pair(ctx, amToken, vToken)
	if err != nil {
		p.logger.Fatal().Err(err).Msg("Pairing failure.")
	}

	p.logger.Info().Msg("Pairing success.")

	return subcommands.ExitSuccess
}
