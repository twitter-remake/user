package config

// Config is the global config for the application
type Config struct {
	AppName     string
	Host        string
	Port        string
	DatabaseURL string
	Environment string

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
	AWSRegion          string
	S3Bucket           string
}

var config *Config

// Init initializes the config package
func Init() {
	config = &Config{
		AppName:            LookupEnv("APP_NAME", "user-service"),
		Host:               LookupEnv("HOST", "0.0.0.0"),
		Port:               LookupEnv("PORT", "9001"),
		DatabaseURL:        LookupEnv("DATABASE_URL", ""),
		Environment:        LookupEnv("ENVIRONMENT", "development"),
		AWSAccessKeyID:     LookupEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: LookupEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSSessionToken:    LookupEnv("AWS_SESSION_TOKEN", ""),
		AWSRegion:          LookupEnv("AWS_REGION", "us-east-1"),
		S3Bucket:           LookupEnv("S3_BUCKET", ""),
	}
}

// GetConfig returns the global config
func GetConfig() *Config { return config }

// AppName returns the app name
func AppName() string { return config.AppName }

// Host returns the host
func Host() string { return config.Host }

// Port returns the port
func Port() string { return config.Port }

// DatabaseURL returns the database url
func DatabaseURL() string { return config.DatabaseURL }

// Environment returns the environment
func Environment() string { return config.Environment }

func AWSAccessKeyID() string     { return config.AWSAccessKeyID }
func AWSSecretAccessKey() string { return config.AWSSecretAccessKey }
func AWSSessionToken() string    { return config.AWSSessionToken }
func AWSRegion() string          { return config.AWSRegion }
func S3Bucket() string           { return config.S3Bucket }
