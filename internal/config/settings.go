package config

import "fmt"

// Settings contains the application config
type Settings struct {
	Environment                    string `yaml:"ENVIRONMENT"`
	Port                           string `yaml:"PORT"`
	GRPCPort                       string `yaml:"GRPC_PORT"`
	UsersAPIGRPCAddr               string `yaml:"USERS_API_GRPC_ADDR"`
	LogLevel                       string `yaml:"LOG_LEVEL"`
	DBUser                         string `yaml:"DB_USER"`
	DBPassword                     string `yaml:"DB_PASSWORD"`
	DBPort                         string `yaml:"DB_PORT"`
	DBHost                         string `yaml:"DB_HOST"`
	DBName                         string `yaml:"DB_NAME"`
	DBMaxOpenConnections           int    `yaml:"DB_MAX_OPEN_CONNECTIONS"`
	DBMaxIdleConnections           int    `yaml:"DB_MAX_IDLE_CONNECTIONS"`
	ServiceName                    string `yaml:"SERVICE_NAME"`
	JwtKeySetURL                   string `yaml:"JWT_KEY_SET_URL"`
	DeploymentBaseURL              string `yaml:"DEPLOYMENT_BASE_URL"`
	TorProxyURL                    string `yaml:"TOR_PROXY_URL"`
	SmartcarClientID               string `yaml:"SMARTCAR_CLIENT_ID"`
	SmartcarClientSecret           string `yaml:"SMARTCAR_CLIENT_SECRET"`
	SmartcarTestMode               bool   `yaml:"SMARTCAR_TEST_MODE"`
	SmartcarWebhookID              string `yaml:"SMARTCAR_WEBHOOK_ID"`
	RedisURL                       string `yaml:"REDIS_URL"`
	RedisPassword                  string `yaml:"REDIS_PASSWORD"`
	RedisTLS                       bool   `yaml:"REDIS_TLS"`
	IngestSmartcarURL              string `yaml:"INGEST_SMARTCAR_URL"`
	IngestSmartcarTopic            string `yaml:"INGEST_SMARTCAR_TOPIC"`
	KafkaBrokers                   string `yaml:"KAFKA_BROKERS"`
	DeviceStatusTopic              string `yaml:"DEVICE_STATUS_TOPIC"`
	PrivacyFenceTopic              string `yaml:"PRIVACY_FENCE_TOPIC"`
	TaskRunNowTopic                string `yaml:"TASK_RUN_NOW_TOPIC"`
	NFTInputTopic                  string `yaml:"NFT_INPUT_TOPIC"`
	NFTOutputTopic                 string `yaml:"NFT_OUTPUT_TOPIC"`
	NFTContractAddr                string `yaml:"NFT_CONTRACT_ADDR"`
	NFTChainID                     int    `yaml:"NFT_CHAIN_ID"`
	TaskStopTopic                  string `yaml:"TASK_STOP_TOPIC"`
	TaskCredentialTopic            string `yaml:"TASK_CREDENTIAL_TOPIC"`
	TaskStatusTopic                string `yaml:"TASK_STATUS_TOPIC"`
	EventsTopic                    string `yaml:"EVENTS_TOPIC"`
	ElasticSearchAppSearchHost     string `yaml:"ELASTIC_SEARCH_APP_SEARCH_HOST"`
	ElasticSearchAppSearchToken    string `yaml:"ELASTIC_SEARCH_APP_SEARCH_TOKEN"`
	ElasticSearchAnalyticsHost     string `yaml:"ELASTIC_SEARCH_ANALYTICS_HOST"`
	ElasticSearchAnalyticsUsername string `yaml:"ELASTIC_SEARCH_ANALYTICS_USERNAME"`
	ElasticSearchAnalyticsPassword string `yaml:"ELASTIC_SEARCH_ANALYTICS_PASSWORD"`
	DeviceDataIndexName            string `yaml:"DEVICE_DATA_INDEX_NAME"`
	AWSRegion                      string `yaml:"AWS_REGION"`
	KMSKeyID                       string `yaml:"KMS_KEY_ID"`
	AutoPiAPIToken                 string `yaml:"AUTO_PI_API_TOKEN"`
	SmartcarManagementToken        string `yaml:"SMARTCAR_MANAGEMENT_TOKEN"`
	CIOSiteID                      string `yaml:"CIO_SITE_ID"`
	CIOApiKey                      string `yaml:"CIO_API_KEY"`
	AWSDocumentsBucketName         string `yaml:"AWS_DOCUMENTS_BUCKET_NAME"`
	NFTS3Bucket                    string `yaml:"NFT_S3_BUCKET"`
	DocumentsAWSAccessKeyID        string `yaml:"DOCUMENTS_AWS_ACCESS_KEY_ID"`
	DocumentsAWSSecretsAccessKey   string `yaml:"DOCUMENTS_AWS_SECRET_ACCESS_KEY"`
	DocumentsAWSEndpoint           string `yaml:"DOCUMENTS_AWS_ENDPOINT"`
}

// GetWriterDSN builds the connection string to the db writer - for now same as reader
func (app *Settings) GetWriterDSN(withSearchPath bool) string {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		app.DBUser,
		app.DBPassword,
		app.DBName,
		app.DBHost,
		app.DBPort,
	)
	if withSearchPath {
		dsn = fmt.Sprintf("%s search_path=%s", dsn, app.DBName) // assumption is schema has same name as dbname
	}
	return dsn
}
