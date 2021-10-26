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

	resp, err := conn.Eval("return init()", []interface{}{})
	if err != nil {
		fmt.Println("Error", err)
		fmt.Println("Code", resp.Code)
		return &SessionManager{}, err
	}

	return &seesManager, nil
}

func (conn *SessionManager) GetUserIDByCookie(sessionCookie string) (userID uint64, err error) {
	// resp, err := conn.TarantoolConn.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{sessionCookie})
	// if err != nil {
	// 	return 0, err
	// }

	resp, err := conn.TarantoolConn.Call("check_session", []interface{}{sessionCookie})
	if err != nil {
		fmt.Println("cannot check session", err)
		return 0, err
	}

	if len(resp.Data) == 0 {
		return 0, errors.New("not exixsts this cookie")
	}

	data := resp.Data[0]
	if data == nil {
		return 0, nil
	}

	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return 0, fmt.Errorf("cannot cast data: %v", sessionDataSlice)
	}

	if len(sessionDataSlice) == 0 {
		return 0, errors.New("not exixsts this cookie")
	}
	if sessionDataSlice[0] == nil {
		return 0, nil
	}

	sessionData, ok := sessionDataSlice[1].(uint64)
	if !ok {
		return 0, fmt.Errorf("cannot cast data: %v", sessionDataSlice[0])
	}

	return sessionData, nil
}

func (conn *SessionManager) NewSessionCookie(sessionCookie string, id uint64) error {
	// resp, err := conn.TarantoolConn.Insert("sessions", []interface{}{"sessionCookie", 3})
	resp, err := conn.TarantoolConn.Call("new_session", []interface{}{sessionCookie, id})
	if err != nil {
		return err
	}

	if len(resp.Data) == 0 {
		return errors.New("this cookie already exists")
	}
	return nil
}

func (conn *SessionManager) DeleteSessionCookie(sessionCookie string) error {
	// _, err := conn.TarantoolConn.Delete("sessions", "primary", []interface{}{sessionCookie})
	resp, err := conn.TarantoolConn.Call("delete_session", []interface{}{sessionCookie})
	if err != nil {
		return err
	}

	if len(resp.Data) == 0 {
		return errors.New("this cookie is not exists")
	}

	fmt.Println(resp)

	return nil
}

func (conn *SessionManager) IsSessionByCookie(sessionCookie string) bool {
	resp, err := conn.TarantoolConn.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{sessionCookie})
	if err != nil {
		return false
	}
	if len(resp.Data) == 0 {
		return false
	}
	return true
}
func (conn *SessionManager) DropCookies() {}
