package session

import (
	"dripapp/internal/dripapp/models"
	"errors"
)

type MockSessionDB struct {
	cookies map[string]uint64
}

func NewSessionDB() *MockSessionDB {
	return &MockSessionDB{make(map[string]uint64)}
}

func (db MockSessionDB) GetSessionByCookie(sessionCookie string) (session models.Session, err error) {
	if len(db.cookies) == 0 {
		return models.Session{}, errors.New("cookies is empty map")
	}

	userID, okCookie := db.cookies[sessionCookie]
	if !okCookie {
		return models.Session{}, errors.New("cookie not found")
	}

	return models.Session{UserID: userID, Cookie: sessionCookie}, nil
}

func (db *MockSessionDB) NewSessionCookie(sessionCookie string, userId uint64) error {
	db.cookies[sessionCookie] = userId
	return nil
}

func (db *MockSessionDB) DeleteSessionCookie(sessionCookie string) error {
	if _, ok := db.cookies[sessionCookie]; !ok {
		return errors.New("cookie not found")
	}

	delete(db.cookies, sessionCookie)
	return nil
}

func (db *MockSessionDB) IsSessionByCookie(sessionCookie string) bool {
	for cookieUser := range db.cookies {
		if cookieUser == sessionCookie {
			return true
		}
	}

	return false
}

func (db *MockSessionDB) DropCookies() {
	db.cookies = make(map[string]uint64)
}
