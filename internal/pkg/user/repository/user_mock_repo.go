package MockDB

import (
	"dripapp/internal/pkg/models"
	"errors"
)

type MockDB struct {
	users       map[uint64]*models.User
	swipedUsers map[uint64][]uint64
	tags        map[uint64]string
}

func NewMockDB() *MockDB {
	return &MockDB{make(map[uint64]*models.User), make(map[uint64][]uint64), make(map[uint64]string)}
}

func (db MockDB) GetUser(email string) (*models.User, error) {
	if len(db.users) == 0 {
		return &models.User{}, errors.New("users is empty map")
	}

	currentUser := &models.User{}
	okUser := false
	for _, value := range db.users {
		if value.Email == email {
			currentUser = value
			okUser = true
		}
	}
	if !okUser {
		return &models.User{}, errors.New("User not found")
	}

	return currentUser, nil
}

func (db MockDB) GetUserByID(userID uint64) (*models.User, error) {
	if user, ok := db.users[userID]; ok {
		return user, nil
	}

	return &models.User{}, errors.New("")
}

func (db *MockDB) CreateUser(logUserData *models.LoginUser) (*models.User, error) {
	newID := uint64(len(db.users) + 1)

	db.users[newID] = models.NewUser(newID, logUserData.Email, logUserData.Password)

	return db.users[newID], nil
}

func (db *MockDB) UpdateUser(newUserData *models.User) (err error) {
	db.users[newUserData.ID] = newUserData

	return nil
}

func (db *MockDB) AddSwipedUsers(currentUserId, swipedUserId uint64) error {
	if len(db.users) == 0 {
		return errors.New("users is empty map")
	}

	db.swipedUsers[currentUserId] = append(db.swipedUsers[currentUserId], swipedUserId)
	return nil
}

func (db MockDB) GetNextUserForSwipe(currentUserId uint64) (*models.User, error) {
	if len(db.users) == 0 {
		return &models.User{}, errors.New("users is empty map")
	}
	if len(db.swipedUsers) == 0 {
		for key, value := range db.users {
			if key != currentUserId {
				return value, nil
			}
		}
		return &models.User{}, errors.New("haven't any other users for swipe")
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

	return &models.User{}, errors.New("haven't any other users for swipe")
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

func (db MockDB) IsSwiped(userID, swipedUserID uint64) bool {
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

func (db *MockDB) CreateUserAndProfile(user *models.User) {
	newID := uint64(len(db.users) + 1)

	user.ID = newID

	db.users[newID] = user
}

func (db *MockDB) DropUsers() {
	db.users = make(map[uint64]*models.User)
}

func (db *MockDB) DropSwipes() {
	db.swipedUsers = make(map[uint64][]uint64)
}

func (db *MockDB) CreateTag(text string) {
	newID := uint64(len(db.tags) + 1)
	db.tags[newID] = text
}

func (db *MockDB) GetTags() map[uint64]string {
	return db.tags
}
