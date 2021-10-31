package configs

import (
	"time"
)

type PostgresConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type TarantoolConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type ServerConfig struct {
	Host     string
	Port     string
	SertFile string
	KeyFile  string
}

type FileStorageConfig struct {
	RootFolder       string
	ProfilePhotoPath string
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

type contextID string

var (
	Postgres PostgresConfig

	Tarantool TarantoolConfig

	Server ServerConfig

	FileStorage FileStorageConfig

	Timeouts TimeoutsConfig

	ForContext contextID
)

func SetConfig() {
	// viper.SetConfigFile("config.json")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if viper.GetBool(`debug`) {
	// 	log.Println("Service RUN on DEBUG mode")
	// }

	Postgres = PostgresConfig{
		// Port:     viper.GetString(`database.port`),
		// Host:     viper.GetString(`database.host`),
		// User:     viper.GetString(`database.user`),
		// Password: viper.GetString(`database.pass`),
		// DBName:   viper.GetString(`database.name`),
		Port:     ":5432",
		Host:     "127.0.0.1",
		User:     "admin",
		Password: "lolkek",
		DBName:   "postgres",
	}

	Tarantool = TarantoolConfig{
		// Port:     viper.GetString(`session.port`),
		// Host:     viper.GetString(`session.host`),
		// User:     viper.GetString(`session.user`),
		// Password: viper.GetString(`session.pass`),
		// DBName:   viper.GetString(`session.name`),
		Port:     ":3301",
		Host:     "127.0.0.1",
		User:     "admin",
		Password: "pass",
		DBName:   "drip",
	}

	Server = ServerConfig{
		// Port:     viper.GetString(`server.port`),
		// Host:     viper.GetString(`server.host`),
		// SertFile: viper.GetString(`server.sertFile`),
		// KeyFile:  viper.GetString(`server.keyFile`),
		Host:     "127.0.0.1",
		Port:     ":8000",
		SertFile: "api.ijia.me.crt",
		KeyFile:  "api.ijia.me.key",
	}

	FileStorage = FileStorageConfig{
		RootFolder:       "media",
		ProfilePhotoPath: "profile_photos",
	}

	Timeouts = TimeoutsConfig{
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		ContextTimeout: time.Second * 2,
	}

	ForContext = "userID"
}
