package config

import "time"

type (
	// Config is main struct for configuration
	Config struct {
		ProjectID            string            `json:",omitempty"`
		HTTPConfig           SectionHTTP       `json:",omitempty" mapstructure:"httpconfig"`
		Core                 SectionCore       `json:",omitempty" mapstructure:"core"`
		GINMode              string            `json:",omitempty"`
		LogService           SectionLogService `json:",omitempty"`
		AppsReferenceID      []string          `json:",omitempty"`
		Datadog              DatadogConfig
		Slack                SectionSlackConfig
		PubSub               SectionPubSub
		GRPC                 GRPCConfig
		HTTP                 HTTPConfig
		PostgreSQLDebug      bool
		GCP                  GCPConfig
		LogLevelDebugCodesOK bool
		DefaultTimeout       time.Duration
		GoProfiler           SectionGoProfiler
		GoCloudProfiler      SectionGoCloudProfiler
		User                 SectionUser
	}

	DatadogConfig struct {
		Host            string
		APMHost         string
		ProfilerEnabled bool
	}

	// SectionCore is a struct for Service Configuration
	SectionCore struct {
		Name          string         `json:",omitempty" mapstructure:"name"`
		Environment   string         `json:",omitempty" mapstructure:"environment"`
		Port          string         `json:",omitempty" mapstructure:"port"`
		LogLevel      string         `json:",omitempty" mapstructure:"loglevel"`
		GRPC          SectionGRPC    `json:",omitempty"`
		DBPostgres    SectionDB      `json:",omitempty" mapstructure:"dbpostgres"`
		Redis         RedisAccount   `json:",omitempty" mapstructure:"redis"`
		Mongo         MongoDBAccount `json:",omitempty" mapstructure:"mongo"`
		Timeout       SectionTimeout `json:",omitempty"`
		AppKey        string         `json:",omitempty"`
		Slack         SectionSlack   `json:",omitempty"`
		MaxProccess   int            `json:",omitempty"`
		Version       string         `json:",omitempty"`
		InternalToken string         `json:",omitempty"`
	}

	SectionGRPC struct {
		Port string `json:",omitempty" mapstructure:"port"`
	}

	SectionSlack struct {
		WebhookURL        string `json:"webhook_url,omitempty" mapstructure:"webhookurl"`
		WebhookChannel    string `json:"webhook_channel,omitempty" mapstructure:"webhookchannel"`
		IsEnableSlack     bool   `json:"is_enable_slack,omitempty" mapstructure:"isenableslack"`
		VerificationToken string `json:"verification_token,omitempty" mapstructure:"verificationtoken"`
		BotToken          string `json:"bot_token,omitempty" mapstructure:"bottoken"`
		ChannelID         string `json:"channel_id,omitempty" mapstructure:"channelid"`
		Icon              string `json:"icon,omitempty" mapstructure:"icon"`
	}

	// SectionHTTP is a config for http request
	SectionHTTP struct {
		Timeout          int  `json:",omitempty" mapstructure:"timeout"`
		DisableKeepAlive bool `json:",omitempty" mapstructure:"disablekeepalive"`
	}

	// SectionTimeout for service timeout in second
	SectionTimeout struct {
		Read  int `json:",omitempty"`
		Write int `json:",omitempty"`
	}

	SectionAWSItem struct {
		AccessKeyID     string `json:",omitempty"`
		SecretAccessKey string `json:",omitempty"`
		DefaultBucket   string `json:",omitempty"`
		Region          string `json:",omitempty"`
		ScoutPrefix     string `json:",omitempty"`
		Host            string `json:",omitempty"`
		Storage         string `json:",omitempty"`
	}

	RedisAccount struct {
		Host                 string   `json:",omitempty" mapstructure:"host"`
		Port                 int      `json:",omitempty" mapstructure:"port"`
		DB                   int      `json:",omitempty" mapstructure:"db"`
		Password             string   `json:",omitempty" mapstructure:"password"`
		PoolSize             int      `json:",omitempty" mapstructure:"poolsize"`
		MinIdleConns         int      `json:",omitempty" mapstructure:"minidleconns"`
		MaxIdle              int      `json:",omitempty" mapstructure:"maxidle"`
		MaxActive            int      `json:",omitempty" mapstructure:"maxactive"`
		MaxConnLifetime      int      `json:",omitempty" mapstructure:"maxconnlifetime"`
		RedisearchIndex      []string `json:",omitempty"`
		RedisMutexExpiryTime int      `json:",omitempty"`
		RedisMutexLockTries  int      `json:",omitempty"`
	}

	SectionDB struct {
		READ  DBAccount `json:",omitempty" mapstructure:"read"`
		WRITE DBAccount `json:",omitempty" mapstructure:"write"`
	}

	// DBAccount is struct for database configuration or database account
	DBAccount struct {
		Username     string `json:",omitempty" mapstructure:"username"`
		Password     string `json:",omitempty" mapstructure:"password"`
		Host         string `json:",omitempty" mapstructure:"host"`
		Port         string `json:",omitempty" mapstructure:"port"`
		Name         string `json:",omitempty" mapstructure:"dbname"`
		Flavor       string `json:",omitempty" mapstructure:"flavor"`
		MaxIdleConns int    `json:",omitempty" mapstructure:"maxidleconns"`
		MaxOpenConns int    `json:",omitempty" mapstructure:"maxopenconns"`
		MaxLifeTime  int    `json:",omitempty" mapstructure:"maxlifetime"`
		Location     string `json:",omitempty" mapstructure:"location"`
		Timeout      string `json:",omitempty" mapstructure:"timeout"`
	}

	// MongoDBAccount holds the account information of MongoDB instance used by IMS
	MongoDBAccount struct {
		Username string `json:",omitempty" mapstructure:"username"`
		Password string `json:",omitempty" mapstructure:"password"`
		Host     string `json:",omitempty" mapstructure:"host"`
		Port     string `json:",omitempty" mapstructure:"port"`
		DBName   string `json:",omitempty" mapstructure:"dbname"`
	}

	SectionLogService struct {
		RequestLoggingEnable bool `json:",omitempty"`
	}

	// option defines configuration option
	option struct {
		configType string
	}

	Option func(*option)

	SectionSlackConfig struct {
		General              SectionSlack
		RollbackNotification SectionSlack
	}

	SectionPubSub struct {
		Publisher  SectionPublisher
		Subscriber SectionSubscriber

		BackoffInitial    time.Duration
		BackoffMax        time.Duration
		BackoffMultiplier float64
	}

	SectionPublisher struct {
		DemoPublisher SectionPubSubDetail
	}

	SectionSubscriber struct {
		DemoSubscriber SectionPubSubDetail
	}

	SectionPubSubDetail struct {
		Topic  string
		Sub    string
		Active bool
	}

	GRPCConfig struct {
		// external grpc
	}

	HTTPConfig struct {
		// http client HTTPItemConfig
		Debug bool
	}

	GRPCItemConfig struct {
		Host                    string
		Port                    int
		MaxRetry                int
		EnableRetry             bool
		TLS                     bool
		Timeout                 int
		ExponentialBackoffRetry int
	}

	HTTPItemConfig struct {
		Host                    string
		Port                    int
		MaxRetry                int
		EnableRetry             bool
		TLS                     bool
		Timeout                 int
		ExponentialBackoffRetry int
	}

	GCPConfig struct {
		ProjectId         string
		StorageBucketName string
	}

	SectionGoProfiler struct {
		Enable bool
		Port   string
	}

	SectionGoCloudProfiler struct {
		Enable        bool
		OpenTelemetry bool
		DebugLogging  bool
	}

	SectionUser struct {
		TTLMinutes time.Duration
	}
)

func (c *Config) IsProduction() bool {
	return c.Core.Environment == "production"
}
