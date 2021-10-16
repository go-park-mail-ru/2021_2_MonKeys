package session

import (
	"errors"
	"fmt"
	"log"

	"github.com/tarantool/go-tarantool"
)

type SessionManager struct {
	TarantoolConn *tarantool.Connection
}

func NewTarantoolConnection() (*SessionManager, error) {
	conn, err := tarantool.Connect("127.0.0.1:3301", tarantool.Opts{
		User: "admin",
		Pass: "pass",
	})

	seesManager := SessionManager{conn}

	if err != nil {
		return &SessionManager{}, err
	} else {
		log.Print("Connection success (mytarantool)")
	}
	return &seesManager, nil
}

func (conn *SessionManager) GetUserIDByCookie(sessionCookie string) (userID uint64, err error) {
	resp, err := conn.TarantoolConn.Select("cookie", "secondary", 0, 1, tarantool.IterEq, []interface{}{sessionCookie})
	if err != nil {
		return 0, err
	}

	fmt.Println(resp.Code)
	if len(resp.Data) == 0 {
		return 0, errors.New("ahahah")
	}
	data := resp.Data[0]
	fmt.Println(data)
	return 0, nil
}

func (conn *SessionManager) NewSessionCookie(sessionCookie string, id uint64) error {
	resp, err := conn.TarantoolConn.Insert("cookie", []interface{}{id, sessionCookie})
	if err != nil {
		return err
	}

	fmt.Println(resp.Code)
	if len(resp.Data) == 0 {
		return errors.New("ahahah")
	}
	return nil
}

func (conn *SessionManager) DeleteSessionCookie(sessionCookie string) error {
	return nil
}
