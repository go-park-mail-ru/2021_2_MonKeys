package main

import (
	"errors"
)

type MockDB struct {
	users       map[uint64]User
	swipedUsers map[uint64][]uint64
}

func NewMockDB() *MockDB {
	return &MockDB{make(map[uint64]User), make(map[uint64][]uint64)}
}

func (db MockDB) getUser(email string) (User, error) {
	if len(db.users) == 0 {
		return User{}, errors.New("users is empty map")
	}

	currentUser := User{}
	okUser := false
	for _, value := range db.users {
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

func (db MockDB) getUserByID(userID uint64) (User, error) {
	if user, ok := db.users[userID]; ok {
		return user, nil
	}

	return User{}, errors.New("")
}

func (db *MockDB) createUser(logUserData LoginUser) (User, error) {
	newID := uint64(len(db.users) + 1)

	db.users[newID] = makeUser(newID, logUserData.Email, logUserData.Password)

	return db.users[newID], nil
}

func (db *MockDB) updateUser(newUserData User) (err error) {
	db.users[newUserData.ID] = newUserData

	return nil
}

func (db *MockDB) addSwipedUsers(currentUserId, swipedUserId uint64) error {
	if len(db.users) == 0 {
		return errors.New("users is empty map")
	}

	db.swipedUsers[currentUserId] = append(db.swipedUsers[currentUserId], swipedUserId)
	return nil
}

func (db MockDB) getNextUserForSwipe(currentUserId uint64) (User, error) {
	if len(db.users) == 0 {
		return User{}, errors.New("users is empty map")
	}
	if len(db.swipedUsers) == 0 {
		for key, value := range db.users {
			if key != currentUserId {
				return value, nil
			}
		}
		return User{}, errors.New("haven't any other users for swipe")
	}

	// find all users swiped by the current user
	var allSwipedUsersForCurrentUser []uint64
	for key, value := range db.swipedUsers {
		if key == currentUserId {
			allSwipedUsersForCurrentUser = value
		}
	}

	// find a user who has not yet been swiped by the current user
	for key, value := range db.users {
		if key == currentUserId {
			continue
		}
		if !existsIn(key, allSwipedUsersForCurrentUser) {
			return value, nil
		}
	}

	return User{}, errors.New("haven't any other users for swipe")
}

func existsIn(value uint64, target []uint64) bool {
	exists := false
	for i := range target {
		if value == target[i] {
			exists = true
		}
	}

	return exists
}

func (db MockDB) isSwiped(userID, swipedUserID uint64) bool {
	swipedUsers, ok := db.swipedUsers[userID]
	if !ok {
		return false
	}

	for _, currentUserID := range swipedUsers {
		if currentUserID == swipedUserID {
			return true
		}
	}

	return false
}

type MockSessionDB struct {
	cookies map[string]uint64
}

func NewSessionDB() *MockSessionDB {
	return &MockSessionDB{make(map[string]uint64)}
}

func (db MockSessionDB) getUserIDByCookie(sessionCookie string) (userID uint64, err error) {
	if len(db.cookies) == 0 {
		return userID, errors.New("cookies is empty map")
	}

	userID, okCookie := db.cookies[sessionCookie]
	if !okCookie {
		return userID, errors.New("cookie not found")
	}

	return userID, nil
}

func (db *MockSessionDB) newSessionCookie(sessionCookie string, userId uint64) error {
	db.cookies[sessionCookie] = userId
	return nil
}

func (db *MockSessionDB) deleteSessionCookie(sessionCookie string) error {
	if _, ok := db.cookies[sessionCookie]; !ok {
		return errors.New("cookie not found")
	}

	delete(db.cookies, sessionCookie)
	return nil
}

func (db MockSessionDB) isSessionByUserID(userID uint64) bool {
	for _, currentUserID := range db.cookies {
		if currentUserID == userID {
			return true
		}
	}

	return false
}
