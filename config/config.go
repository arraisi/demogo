package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *Config

// Init initializes `config` from the default config file.
// use `WithConfigFile` to specify the location of the config file
func Init(opts ...Option) (*Config, error) {
	configFile := ".env"
	viper.AutomaticEnv()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return &Config{}, err
		}
	}
	config = new(Config)
	// GIN
	config.GINMode = viper.GetString("GIN_MODE")

	// HTTP Config
	config.HTTPConfig.Timeout = viper.GetInt("HTTP_TIMEOUT_SECOND")
	config.HTTPConfig.DisableKeepAlive = viper.GetBool("HTTP_DISABLE_KEEP_ALIVE")

	config.Core.Redis.Host = viper.GetString("REDIS_HOST")
	config.Core.Redis.Port = viper.GetInt("REDIS_PORT")
	config.Core.Redis.DB = viper.GetInt("REDIS_DB")
	config.Core.Redis.Password = viper.GetString("REDIS_PASSWORD")
	config.Core.Redis.PoolSize = viper.GetInt("REDIS_POOL_SIZE")
	config.Core.Redis.MinIdleConns = viper.GetInt("REDIS_MIN_IDLE_CONNS")

	config.Core.DBPostgres.READ.Username = viper.GetString("POSTGRES_DATABASE_READ_USERNAME")
	config.Core.DBPostgres.READ.Password = viper.GetString("POSTGRES_DATABASE_READ_PASSWORD")
	config.Core.DBPostgres.READ.Host = viper.GetString("POSTGRES_DATABASE_READ_HOST")
	config.Core.DBPostgres.READ.Port = viper.GetString("POSTGRES_DATABASE_READ_PORT")
	config.Core.DBPostgres.READ.Name = viper.GetString("POSTGRES_DATABASE_READ_NAME")
	config.Core.DBPostgres.READ.MaxIdleConns = viper.GetInt("POSTGRES_DATABASE_READ_MAXIDLECONNS")
	config.Core.DBPostgres.READ.MaxOpenConns = viper.GetInt("POSTGRES_DATABASE_READ_MAXOPENCONNS")
	config.Core.DBPostgres.READ.MaxLifeTime = viper.GetInt("POSTGRES_DATABASE_READ_MAXLIFETIME")
	config.Core.DBPostgres.READ.Timeout = viper.GetString("POSTGRES_DATABASE_READ_TIMEOUT")

	config.Core.DBPostgres.WRITE.Username = viper.GetString("POSTGRES_DATABASE_WRITE_USERNAME")
	config.Core.DBPostgres.WRITE.Password = viper.GetString("POSTGRES_DATABASE_WRITE_PASSWORD")
	config.Core.DBPostgres.WRITE.Host = viper.GetString("POSTGRES_DATABASE_WRITE_HOST")
	config.Core.DBPostgres.WRITE.Port = viper.GetString("POSTGRES_DATABASE_WRITE_PORT")
	config.Core.DBPostgres.WRITE.Name = viper.GetString("POSTGRES_DATABASE_WRITE_NAME")
	config.Core.DBPostgres.WRITE.MaxIdleConns = viper.GetInt("POSTGRES_DATABASE_WRITE_MAXIDLECONNS")
	config.Core.DBPostgres.WRITE.MaxOpenConns = viper.GetInt("POSTGRES_DATABASE_WRITE_MAXOPENCONNS")
	config.Core.DBPostgres.WRITE.MaxLifeTime = viper.GetInt("POSTGRES_DATABASE_WRITE_MAXLIFETIME")
	config.Core.DBPostgres.WRITE.Timeout = viper.GetString("POSTGRES_DATABASE_WRITE_TIMEOUT")

	config.Core.Name = viper.GetString("SERVICE_NAME")
	config.Core.Environment = viper.GetString("SERVICE_ENVIRONMENT")
	config.Core.Port = viper.GetString("SERVICE_PORT")
	config.Core.LogLevel = viper.GetString("SERVICE_LOG_LEVEL")
	config.Core.Version = viper.GetString("SERVICE_VERSION")

	config.Datadog.Host = viper.GetString("DATADOG_HOST")
	config.Datadog.APMHost = viper.GetString("DATADOG_APM_HOST")
	config.Datadog.ProfilerEnabled = viper.GetBool("DD_PROFILER_ENABLED")

	config.Core.GRPC.Port = viper.GetString("SERVICE_GRPC_PORT")

	config.Core.Mongo.Host = viper.GetString("MONGO_HOST")
	config.Core.Mongo.Port = viper.GetString("MONGO_PORT")
	config.Core.Mongo.Username = viper.GetString("MONGO_USERNAME")
	config.Core.Mongo.Password = viper.GetString("MONGO_PASSWORD")
	config.Core.Mongo.DBName = viper.GetString("MONGO_DB_NAME")

	config.LogService.RequestLoggingEnable = viper.GetBool("SERVICE_REQUEST_LOGGING_ENABLED")

	config.AppsReferenceID = strings.Split(viper.GetString("SERVICE_APPS_REFERENCE_ID"), ",")
	config.ProjectID = viper.GetString("SERVICE_GCP_PROJECT_ID")

	config.PostgreSQLDebug = viper.GetBool("POSTGRESQL_DEBUG_MODE")

	config.GCP.ProjectId = viper.GetString("GCP_PROJECT_ID")
	config.GCP.StorageBucketName = viper.GetString("GCP_BUCKET_NAME")

	config.LogLevelDebugCodesOK = viper.GetBool("LOG_LEVEL_DEBUG_CODES_OK")

	config.DefaultTimeout = viper.GetDuration("DEFAULT_TIMEOUT_SECONDS")

	config.GoProfiler.Enable = viper.GetBool("GO_PROFILER_ENABLED")
	config.GoProfiler.Port = viper.GetString("GO_PROFILER_PORT")

	config.GoCloudProfiler.Enable = viper.GetBool("GO_CLOUD_PROFILER_ENABLED")
	config.GoCloudProfiler.OpenTelemetry = viper.GetBool("GO_CLOUD_PROFILER_OPEN_TELEMETRY_ENABLED")
	config.GoCloudProfiler.DebugLogging = viper.GetBool("GO_CLOUD_PROFILER_DEBUG_LOGGING_ENABLED")

	config.User.TTLMinutes = viper.GetDuration("USER_TTL_MINUTES")

	if err := viper.Unmarshal(config); err != nil {
		return &Config{}, err
	}

	return config, nil
}
