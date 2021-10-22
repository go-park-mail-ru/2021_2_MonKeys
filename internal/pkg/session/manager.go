package session

import (
	"dripapp/configs"
	"errors"
	"fmt"
	"log"

	"github.com/tarantool/go-tarantool"
)

const success = "Connection success (drip_tarantool) on: "

type SessionManager struct {
	TarantoolConn *tarantool.Connection
}

func NewTarantoolConnection(tntConfig configs.TarantoolConfig) (*SessionManager, error) {
	addrPort := fmt.Sprintf("%s%s", tntConfig.Host, tntConfig.Port)
	conn, err := tarantool.Connect(addrPort, tarantool.Opts{
		User: "guest",
		// User: tntConfig.User,
		// Pass: tntConfig.Password,
	})

	seesManager := SessionManager{conn}

	if err != nil {
		return &SessionManager{}, err
	} else {
		log.Printf("%s%s", success, addrPort)
	}

	return &seesManager, nil
}

func (conn *SessionManager) GetUserIDByCookie(sessionCookie string) (userID uint64, err error) {
	resp, err := conn.TarantoolConn.Select("sessions", "secondary", 0, 1, tarantool.IterEq, []interface{}{sessionCookie})
	if err != nil {
		return 0, err
	}

	if len(resp.Data) == 0 {
		return 0, errors.New("not exixsts this cookie")
	}
	data := resp.Data[0]
	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return 0, fmt.Errorf("cannot cast data: %v", sessionDataSlice)
	}
	return sessionDataSlice[0].(uint64), nil
}

func (conn *SessionManager) NewSessionCookie(sessionCookie string, id uint64) error {
	resp, err := conn.TarantoolConn.Insert("sessions", []interface{}{id, sessionCookie})
	if err != nil {
		return err
	}

	if len(resp.Data) == 0 {
		return errors.New("this cookie already exists")
	}
	return nil
}

func (conn *SessionManager) DeleteSessionCookie(sessionCookie string) error {
	_, err := conn.TarantoolConn.Delete("sessions", "secondary", []interface{}{sessionCookie})
	return err
}

func (conn *SessionManager) IsSessionByUserID(userID uint64) bool {
	resp, err := conn.TarantoolConn.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(userID)})
	if err != nil {
		return false
	}
	if len(resp.Data) == 0 {
		return false
	}
	return true
}
func (conn *SessionManager) DropCookies() {}
