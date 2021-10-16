package session

import "errors"

type MockSessionDB struct {
	cookies map[string]uint64
}

func NewSessionDB() *MockSessionDB {
	return &MockSessionDB{make(map[string]uint64)}
}

func (db MockSessionDB) GetUserIDByCookie(sessionCookie string) (userID uint64, err error) {
	if len(db.cookies) == 0 {
		return userID, errors.New("cookies is empty map")
	}

	userID, okCookie := db.cookies[sessionCookie]
	if !okCookie {
		return userID, errors.New("cookie not found")
	}

	return userID, nil
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

func (db *MockSessionDB) IsSessionByUserID(userID uint64) bool {
	for _, currentUserID := range db.cookies {
		if currentUserID == userID {
			return true
		}
	}

	return false
}

func (db *MockSessionDB) DropCookies() {
	db.cookies = make(map[string]uint64)
}
