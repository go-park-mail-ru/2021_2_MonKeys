package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
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

type timeouts struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

var Postgres PostgresConfig

var Tarantool TarantoolConfig

var Server ServerConfig

var Timeouts timeouts

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	// Postgres = PostgresConfig{
	// 	User:     os.Getenv("PostgresUser"),
	// 	Password: os.Getenv("PostgresPassword"),
	// 	Port:     os.Getenv("PostgresPort"),
	// 	Host:     os.Getenv("PostgresHost"),
	// 	DBName:   os.Getenv("PostgresDBName"),
	// }

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
		SertFile: viper.GetString(`server.sertFile`),
		KeyFile:  viper.GetString(`server.keyFile`),
	}

	Timeouts = timeouts{
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		ContextTimeout: time.Second * 2,
	}
}
