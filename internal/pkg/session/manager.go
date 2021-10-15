package session

import (
	"log"

	"github.com/tarantool/go-tarantool"
)

type SessionManager struct {
	TarantoolConn *tarantool.Connection
}

func NewTarantoolConnection() *SessionManager {
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "pass",
	})

	seesManager := SessionManager{conn}

	if err != nil {
		log.Fatalf("Connection refused")
		panic(err)
	} else {
		log.Print("Connection success (mytarantool)")
	}
	return &seesManager
}
