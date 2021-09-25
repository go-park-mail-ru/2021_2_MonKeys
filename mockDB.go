package main

import "errors"

var (
	users   = make(map[uint64]User)
	cookies = make(map[string]uint64)
)

type MockDB struct {
	//DB int
}

func (MockDB) getUserModel(email string) (User, error) {
	if len(users) == 0 {
		return User{}, errors.New("users is empty map")
	}

	currentUser := User{}
	okUser := false
	for _, value := range users {
		if value.Email == email {
			currentUser = value
			okUser = true
		}
	}
	if !okUser {
		return User{}, errors.New("User not found")
	}

	return currentUser, nil
}

type MockSessionDB struct {
}

func (MockSessionDB) getUserByCookie(sessionCookie string) (User, error) {
	if len(cookies) == 0 {
		return User{}, errors.New("cookies is empty map")
	}

	currentUserId, okCookie := cookies[sessionCookie]
	if !okCookie {
		return User{}, errors.New("cookie not found")
	}

	currentUser, okUser := users[currentUserId]
	if !okUser {
		return User{}, errors.New("user not found")
	}

	return currentUser, nil
}

func (MockSessionDB) newSessionCookie(sessionCookie string, userId uint64) error {
	cookies[sessionCookie] = userId
	return nil
}

func (MockSessionDB) deleteSessionCookie(sessionCookie string) error {
	delete(cookies, sessionCookie)
	return nil
}
