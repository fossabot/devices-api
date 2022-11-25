package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/DIMO-Network/devices-api/docs"
	"github.com/DIMO-Network/devices-api/internal/config"
	"github.com/DIMO-Network/devices-api/internal/database"
	es "github.com/DIMO-Network/devices-api/internal/elasticsearch"
	"github.com/DIMO-Network/devices-api/internal/kafka"
	"github.com/DIMO-Network/devices-api/internal/services"
	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/burdiyan/kafkautil"
	"github.com/customerio/go-customerio/v3"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	_ "go.uber.org/automaxprocs"
)

// @title                      DIMO Devices API
// @version                    1.0
// @BasePath                   /v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	gitSha1 := os.Getenv("GIT_SHA1")
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Str("git-sha1", gitSha1).
		Logger()

	config.SetupMachineryLogging(&logger)

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load settings")
	}
	level, err := zerolog.ParseLevel(settings.LogLevel)
	if err != nil {
		logger.Fatal().Err(err).Msgf("could not parse LOG_LEVEL: %s", settings.LogLevel)
	}
	zerolog.SetGlobalLevel(level)

	deps := newDependencyContainer(&settings, logger)

	pdb := database.NewDbConnectionFromSettings(ctx, &settings, true)
	// check db ready, this is not ideal btw, the db connection handler would be nicer if it did this.
	totalTime := 0
	for !pdb.IsReady() {
		if totalTime > 30 {
			logger.Fatal().Msg("could not connect to postgres after 30 seconds")
		}
		time.Sleep(time.Second)
		totalTime++
	}

	esInstance, err := es.NewElasticSearch(settings, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Couldn't instantiate Elasticsearch client.")
	}

	nhtsaSvc := services.NewNHTSAService()
	ddSvc := services.NewDeviceDefinitionService(pdb.DBS, &logger, nhtsaSvc, &settings)
	// todo: use flag or other package to handle args
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	switch arg {
	case "migrate":
		command := "up"
		if len(os.Args) > 2 {
			command = os.Args[2]
			if command == "down-to" || command == "up-to" {
				command = command + " " + os.Args[3]
			}
		}
		migrateDatabase(logger, &settings, command)
	case "generate-events":
		eventService := services.NewEventService(&logger, &settings, deps.getKafkaProducer())
		nhtsaSvc := services.NewNHTSAService()
		ddSvc := services.NewDeviceDefinitionService(pdb.DBS, &logger, nhtsaSvc, &settings)
		generateEvents(logger, pdb, eventService, ddSvc)
	case "set-command-compat":
		if err := setCommandCompatibility(ctx, &settings, pdb, ddSvc); err != nil {
			logger.Fatal().Err(err).Msg("Failed during command compatibility fill.")
		}
		logger.Info().Msg("Finished setting command compatibility.")
	case "remake-smartcar-topic":
		err = remakeSmartcarTopic(ctx, pdb, deps.getKafkaProducer(), ddSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running Smartcar Kafka re-registration")
		}
	case "remake-autopi-topic":
		err = remakeAutoPiTopic(ctx, pdb, deps.getKafkaProducer(), ddSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running AutoPi Kafka re-registration")
		}
	case "remake-fence-topic":
		err = remakeFenceTopic(&settings, pdb, deps.getKafkaProducer())
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running Smartcar Kafka re-registration")
		}
	case "remake-dd-topics":
		err = remakeDeviceDefinitionTopics(ctx, &settings, pdb, deps.getKafkaProducer(), &logger, ddSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error recreating device definition KTables.")
		}
	case "populate-es-dd-data":
		err = populateESDDData(ctx, &settings, esInstance, pdb, &logger, ddSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running elastic search dd update")
		}
	case "populate-es-region-data":
		err = populateESRegionData(ctx, &settings, esInstance, pdb, &logger, ddSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running elastic search region update")
		}
	case "populate-usa-powertrain":
		logger.Info().Msg("Populating USA powertrain data from VINs")
		nhtsaSvc := services.NewNHTSAService()
		err := populateUSAPowertrain(ctx, &logger, pdb, nhtsaSvc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error filling in powertrain data.")
		}
	case "stop-task-by-key":
		if len(os.Args[1:]) != 2 {
			logger.Fatal().Msgf("Expected an argument, the task key.")
		}
		taskKey := os.Args[2]
		logger.Info().Msgf("Stopping task %s", taskKey)
		err := stopTaskByKey(&settings, taskKey, deps.getKafkaProducer())
		if err != nil {
			logger.Fatal().Err(err).Msg("Error stopping task.")
		}
	case "start-smartcar-from-refresh":
		if len(os.Args[1:]) != 2 {
			logger.Fatal().Msgf("Expected an argument, the device ID.")
		}
		userDeviceID := os.Args[2]
		logger.Info().Msgf("Trying to start Smartcar task for %s.", userDeviceID)
		var cipher shared.Cipher
		if settings.Environment == "dev" || settings.Environment == "prod" {
			cipher = createKMS(&settings, &logger)
		} else {
			logger.Warn().Msg("Using ROT13 encrypter. Only use this for testing!")
			cipher = new(shared.ROT13Cipher)
		}
		scClient := services.NewSmartcarClient(&settings)
		scTask := services.NewSmartcarTaskService(&settings, deps.getKafkaProducer())
		if err := startSmartcarFromRefresh(ctx, &logger, &settings, pdb, cipher, userDeviceID, scClient, scTask, ddSvc); err != nil {
			logger.Fatal().Err(err).Msg("Error starting Smartcar task.")
		}
		logger.Info().Msgf("Successfully started Smartcar task for %s.", userDeviceID)
	case "drivly-sync-data":
		logger.Info().Msgf("Sync driv.ly")
		err = loadUserDeviceDrively(ctx, &logger, &settings, pdb)
		if err != nil {
			logger.Fatal().Err(err).Msg("error trying to sync driv.ly")
		}
	case "blackbook-sync-data":
		logger.Info().Msgf("Sync blackbook")
		err = loadUserDeviceBlackbook(ctx, &logger, &settings, pdb)
		if err != nil {
			logger.Fatal().Err(err).Msg("error trying to sync blackbook")
		}
	default:
		startMonitoringServer(logger)
		eventService := services.NewEventService(&logger, &settings, deps.getKafkaProducer())
		startDeviceStatusConsumer(logger, &settings, pdb, eventService)
		startCredentialConsumer(logger, &settings, pdb)
		startTaskStatusConsumer(logger, &settings, pdb)
		startWebAPI(logger, &settings, pdb, eventService, deps.getKafkaProducer(), deps.getS3ServiceClient(ctx), deps.getS3NFTServiceClient(ctx))
	}
}

func createKafkaProducer(settings *config.Settings) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner
	p, err := sarama.NewSyncProducer(strings.Split(settings.KafkaBrokers, ","), config)
	if err != nil {
		return nil, fmt.Errorf("failed to construct producer with broker list %s: %w", settings.KafkaBrokers, err)
	}
	return p, nil
}

func createKMS(settings *config.Settings, logger *zerolog.Logger) shared.Cipher {
	// Need AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to be set.
	// TODO(elffjs): Can we let the SDK grab the region too?
	awscfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(settings.AWSRegion))
	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't create AWS config.")
	}

	return &shared.KMSCipher{
		KeyID:  settings.KMSKeyID,
		Client: kms.NewFromConfig(awscfg),
	}
}

func changeLogLevel(c *fiber.Ctx) error {
	payload := struct {
		LogLevel string `json:"logLevel"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	level, err := zerolog.ParseLevel(payload.LogLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	return c.Status(fiber.StatusOK).SendString("log level set to: " + level.String())
}

func startDeviceStatusConsumer(logger zerolog.Logger, settings *config.Settings, pdb database.DbStore, eventService services.EventService) {
	clusterConfig := sarama.NewConfig()
	clusterConfig.Version = sarama.V2_8_1_0
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	cfg := &kafka.Config{
		ClusterConfig:   clusterConfig,
		BrokerAddresses: strings.Split(settings.KafkaBrokers, ","),
		Topic:           settings.DeviceStatusTopic,
		GroupID:         "user-devices",
		MaxInFlight:     int64(5),
	}
	consumer, err := kafka.NewConsumer(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not start device status update consumer")
	}
	nhtsaSvc := services.NewNHTSAService()
	ddSvc := services.NewDeviceDefinitionService(pdb.DBS, &logger, nhtsaSvc, settings)
	ingestSvc := services.NewDeviceStatusIngestService(pdb.DBS, &logger, eventService, ddSvc)
	consumer.Start(context.Background(), ingestSvc.ProcessDeviceStatusMessages)

	logger.Info().Msg("Device status update consumer started")
}

func startCredentialConsumer(logger zerolog.Logger, settings *config.Settings, pdb database.DbStore) {
	clusterConfig := sarama.NewConfig()
	clusterConfig.Version = sarama.V2_8_1_0
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	cfg := &kafka.Config{
		ClusterConfig:   clusterConfig,
		BrokerAddresses: strings.Split(settings.KafkaBrokers, ","),
		Topic:           settings.TaskCredentialTopic,
		GroupID:         "user-devicesYY",
		MaxInFlight:     int64(5),
	}
	consumer, err := kafka.NewConsumer(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not start credential update consumer")
	}
	credService := services.NewCredentialListener(pdb.DBS, &logger)
	consumer.Start(context.Background(), credService.ProcessCredentialsMessages)

	logger.Info().Msg("Credential update consumer started")
}

func startTaskStatusConsumer(logger zerolog.Logger, settings *config.Settings, pdb database.DbStore) {
	clusterConfig := sarama.NewConfig()
	clusterConfig.Version = sarama.V2_8_1_0
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	cfg := &kafka.Config{
		ClusterConfig:   clusterConfig,
		BrokerAddresses: strings.Split(settings.KafkaBrokers, ","),
		Topic:           settings.TaskStatusTopic,
		GroupID:         "user-devices",
		MaxInFlight:     int64(5),
	}
	consumer, err := kafka.NewConsumer(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not start credential update consumer")
	}
	cio := customerio.NewTrackClient(
		settings.CIOSiteID,
		settings.CIOApiKey,
		customerio.WithRegion(customerio.RegionUS),
	)

	nhtsaSvc := services.NewNHTSAService()
	ddSvc := services.NewDeviceDefinitionService(pdb.DBS, &logger, nhtsaSvc, settings)

	taskStatusService := services.NewTaskStatusListener(pdb.DBS, &logger, cio, ddSvc)
	consumer.Start(context.Background(), taskStatusService.ProcessTaskUpdates)

	logger.Info().Msg("Task status consumer started")
}

func startMonitoringServer(logger zerolog.Logger) {
	monApp := fiber.New(fiber.Config{DisableStartupMessage: true})

	monApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	monApp.Put("/loglevel", changeLogLevel)

	go func() {
		// TODO(elffjs): Make the port a setting.
		if err := monApp.Listen(":8888"); err != nil {
			logger.Fatal().Err(err).Str("port", "8888").Msg("Failed to start monitoring web server.")
		}
	}()

	logger.Info().Str("port", "8888").Msg("Started monitoring web server.")
}

// dependencyContainer way to hold different dependencies we need for our app. We could put all our deps and follow this pattern for everything.
type dependencyContainer struct {
	kafkaProducer      sarama.SyncProducer
	settings           *config.Settings
	logger             *zerolog.Logger
	s3ServiceClient    *s3.Client
	s3NFTServiceClient *s3.Client
}

func newDependencyContainer(settings *config.Settings, logger zerolog.Logger) dependencyContainer {
	return dependencyContainer{
		settings: settings,
		logger:   &logger,
	}
}

// getKafkaProducer instantiates a new kafka producer if not already set in our container and returns
func (dc *dependencyContainer) getKafkaProducer() sarama.SyncProducer {
	if dc.kafkaProducer == nil {
		p, err := createKafkaProducer(dc.settings)
		if err != nil {
			dc.logger.Fatal().Err(err).Msg("Could not initialize Kafka producer, terminating")
		}
		dc.kafkaProducer = p
	}
	return dc.kafkaProducer
}

// getS3ServiceClient instantiates a new default config and then a new s3 services client if not already set. Takes context in, although it could likely use a context from container passed in on instantiation
func (dc *dependencyContainer) getS3ServiceClient(ctx context.Context) *s3.Client {
	if dc.s3ServiceClient == nil {

		cfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(dc.settings.AWSRegion),
			// Comment the below out if not using localhost
			awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {

					if dc.settings.Environment == "local" {
						return aws.Endpoint{PartitionID: "aws", URL: dc.settings.DocumentsAWSEndpoint, SigningRegion: dc.settings.AWSRegion}, nil // The SigningRegion key was what's was missing! D'oh.
					}

					// returning EndpointNotFoundError will allow the service to fallback to its default resolution
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				})))

		if err != nil {
			dc.logger.Fatal().Err(err).Msg("Could not load aws config, terminating")
		}

		dc.s3ServiceClient = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.Region = dc.settings.AWSRegion
			o.Credentials = credentials.NewStaticCredentialsProvider(dc.settings.DocumentsAWSAccessKeyID, dc.settings.DocumentsAWSSecretsAccessKey, "")
		})
	}
	return dc.s3ServiceClient
}

func (dc *dependencyContainer) getS3NFTServiceClient(ctx context.Context) *s3.Client {
	if dc.s3NFTServiceClient == nil {

		cfg, err := awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(dc.settings.AWSRegion),
			// Comment the below out if not using localhost
			awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {

					if dc.settings.Environment == "local" {
						return aws.Endpoint{PartitionID: "aws", URL: dc.settings.DocumentsAWSEndpoint, SigningRegion: dc.settings.AWSRegion}, nil // The SigningRegion key was what's was missing! D'oh.
					}

					// returning EndpointNotFoundError will allow the service to fallback to its default resolution
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				})))

		if err != nil {
			dc.logger.Fatal().Err(err).Msg("Could not load aws config, terminating")
		}

		dc.s3NFTServiceClient = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.Region = dc.settings.AWSRegion
			o.Credentials = credentials.NewStaticCredentialsProvider(dc.settings.NFTAWSAccessKeyID, dc.settings.NFTAWSSecretsAccessKey, "")
		})
	}
	return dc.s3NFTServiceClient
}
