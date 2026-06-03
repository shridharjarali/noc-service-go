package config

// Config holds all external service URLs and tuning parameters.
// Values are populated from configs/app.yaml or environment variables.
type Config struct {
	// Server
	ServerPort    int
	ServerContext string

	// Database
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSchema   string

	// Kafka
	KafkaBrokers      []string
	KafkaConsumerGroup string

	// Kafka persister topics
	SaveTopic           string
	UpdateTopic         string
	UpdateWorkflowTopic string

	// ID-gen
	IdGenHost              string
	IdGenPath              string
	ApplicationNoIdgenName string

	// Pagination defaults
	DefaultLimit   int
	DefaultOffset  int
	MaxSearchLimit int

	// Fuzzy search
	IsFuzzyEnabled bool

	// User service
	UserHost           string
	UserContextPath    string
	UserSearchEndpoint string

	// Workflow
	WfHost                      string
	WfTransitionPath            string
	WfBusinessServiceSearchPath string
	WfProcessPath               string

	// MDMS
	MdmsHost     string
	MdmsEndPoint string

	// BPA
	BpaHost           string
	BpaContextPath    string
	BpaSearchEndpoint string

	// SMS / Notification
	SmsNotifTopic              string
	IsSMSEnabled               bool
	LocalizationHost           string
	LocalizationContextPath    string
	LocalizationSearchEndpoint string
	IsLocalizationStateLevel   bool

	// NOC-specific
	NocOfflineDocRequired bool
}
