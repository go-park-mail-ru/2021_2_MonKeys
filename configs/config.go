package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	DEBUG   = 1
	INFO    = 2
	WARNING = 3
	ERROR   = 4
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
	CertFile string
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

type logLevel int

type contextUserID string

type contextUser string

var (
	Postgres PostgresConfig

	Tarantool TarantoolConfig

	Server ServerConfig

	FileStorage FileStorageConfig

	Timeouts TimeoutsConfig

	LogLevel logLevel

	ContextUserID contextUserID

	ContextUser contextUser
)

func SetConfig() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	logLevelStr := viper.GetString(`log_level`)

	log.Printf("Service RUN on %s mode", logLevelStr)

	switch logLevelStr {
	case "DEBUG":
		LogLevel = DEBUG
	case "INFO":
		LogLevel = INFO
	case "WARNING":
		LogLevel = WARNING
	case "ERROR":
		LogLevel = ERROR
	}

	Postgres = PostgresConfig{
		Port:     viper.GetString(`database.port`),
		Host:     viper.GetString(`database.host`),
		User:     viper.GetString(`database.user`),
		Password: viper.GetString(`database.pass`),
		DBName:   viper.GetString(`database.name`),
	}

	Tarantool = TarantoolConfig{
		Port:     viper.GetString(`session.port`),
		Host:     viper.GetString(`session.host`),
		User:     viper.GetString(`session.user`),
		Password: viper.GetString(`session.pass`),
		DBName:   viper.GetString(`session.name`),
	}

	Server = ServerConfig{
		Port:     viper.GetString(`server.port`),
		Host:     viper.GetString(`server.host`),
		CertFile: viper.GetString(`server.certFile`),
		KeyFile:  viper.GetString(`server.keyFile`),
	}

	FileStorage = FileStorageConfig{
		RootFolder:       viper.GetString(`file_storage.root_folder`),
		ProfilePhotoPath: viper.GetString(`file_storage.profile_photo_path`),
	}

	Timeouts = TimeoutsConfig{
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		ContextTimeout: time.Second * 2,
	}

	ContextUserID = "userID"

	ContextUser = "user"
}
