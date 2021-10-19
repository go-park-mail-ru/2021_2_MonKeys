package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type postgresConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type tarantoolConfig struct {
	Size      int
	Network   string
	Address   string
	Password  string
	SecretKey []byte
}

type serverConfig struct {
	Host     string
	Port     string
	SertFile string
	KeyFile  string
}

var Postgres postgresConfig

var Tarantool tarantoolConfig

var Server serverConfig

var Timeouts timeouts

type timeouts struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	// Postgres = postgresConfig{
	// 	User:     os.Getenv("PostgresUser"),
	// 	Password: os.Getenv("PostgresPassword"),
	// 	Port:     os.Getenv("PostgresPort"),
	// 	Host:     os.Getenv("PostgresHost"),
	// 	DBName:   os.Getenv("PostgresDBName"),
	// }

	// Tarantool = tarantoolConfig{
	// 	Size:      10,
	// 	Network:   "tcp",
	// 	Address:   os.Getenv("RedisAddress"),
	// 	Password:  os.Getenv("RedisPassword"),
	// 	SecretKey: []byte(os.Getenv("SESSION_KEY")),
	// }

	Server = serverConfig{
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
