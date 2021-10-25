package repository

import (
	"context"
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

func (newDB *MockDB) MockDB() {
	newDB.CreateUserAndProfile(context.TODO(), &models.User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "lol1@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		Imgs:        []string{"/img/Yachty-tout.jpg"},
		Tags:        []string{"soccer", "anime"},
	})
	newDB.CreateUserAndProfile(context.TODO(), &models.User{
		ID:          2,
		Name:        "Mikhail2",
		Email:       "lol2@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		Imgs:        []string{"/img/Yachty-tout.jpg"},
		Tags:        []string{"soccer", "anime"},
	})
	newDB.CreateUserAndProfile(context.TODO(), &models.User{
		ID:          3,
		Name:        "Mikhail3",
		Email:       "lol3@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		Imgs:        []string{"/img/Yachty-tout.jpg"},
		Tags:        []string{"soccer", "anime"},
	})
	newDB.CreateUserAndProfile(context.TODO(), &models.User{
		ID:          4,
		Name:        "Mikhail4",
		Email:       "lol4@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		Imgs:        []string{"/img/Yachty-tout.jpg"},
		Tags:        []string{"soccer", "anime"},
	})
	newDB.CreateTag(context.TODO(), "anime")
	newDB.CreateTag(context.TODO(), "netflix")
	newDB.CreateTag(context.TODO(), "games")
	newDB.CreateTag(context.TODO(), "walk")
	newDB.CreateTag(context.TODO(), "JS")
	newDB.CreateTag(context.TODO(), "baumanka")
	newDB.CreateTag(context.TODO(), "music")
	newDB.CreateTag(context.TODO(), "sport")
}

func (db *MockDB) GetUser(ctx context.Context, email string) (*models.User, error) {
	if len(db.users) == 0 {
		return &models.User{}, errors.New("users is empty map")
	}

	currentUser := models.User{}
	okUser := false
	for _, value := range db.users {
		if value.Email == email {
			currentUser = *value
			okUser = true
		}
	}
	if !okUser {
		return &models.User{}, errors.New("User not found")
	}

	return &currentUser, nil
}

func (db *MockDB) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	if user, ok := db.users[userID]; ok {
		return user, nil
	}

	return &models.User{}, errors.New("")
}

func (db *MockDB) CreateUser(ctx context.Context, logUserData *models.LoginUser) (*models.User, error) {
	newID := uint64(len(db.users) + 1)

	db.users[newID] = models.NewUser(newID, logUserData.Email, logUserData.Password)

	return db.users[newID], nil
}

func (db *MockDB) UpdateUser(ctx context.Context, newUserData *models.User) (err error) {
	db.users[newUserData.ID] = newUserData

	return nil
}

func (db *MockDB) AddSwipedUsers(ctx context.Context, currentUserId, swipedUserId uint64, type_name string) error {
	if len(db.users) == 0 {
		return errors.New("users is empty map")
	}

	db.swipedUsers[currentUserId] = append(db.swipedUsers[currentUserId], swipedUserId)
	return nil
}

func (db *MockDB) GetNextUserForSwipe(ctx context.Context, currentUserId uint64) (models.User, error) {
	if len(db.users) == 0 {
		return models.User{}, errors.New("users is empty map")
	}
	if len(db.swipedUsers) == 0 {
		for key, value := range db.users {
			if key != currentUserId {
				return *value, nil
			}
		}
		return models.User{}, errors.New("haven't any other users for swipe")
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
			return *value, nil
		}
	}

	return models.User{}, errors.New("haven't any other users for swipe")
}

func existsIn(value uint64, target []uint64) bool {
	for i := range target {
		if value == target[i] {
			return true
		}
	}

	return false
}

func (db *MockDB) IsSwiped(ctx context.Context, userID, swipedUserID uint64) (bool, error) {
	swipedUsers, ok := db.swipedUsers[userID]
	if !ok {
		return false, nil
	}

	for _, currentUserID := range swipedUsers {
		if currentUserID == swipedUserID {
			return true, nil
		}
	}

	return false, nil
}

func (db *MockDB) CreateUserAndProfile(ctx context.Context, user *models.User) (models.User, error) {
	newID := uint64(len(db.users) + 1)

	user.ID = newID

	db.users[newID] = user

	return *user, nil
}

func (db *MockDB) DropUsers(ctx context.Context) error {
	db.users = make(map[uint64]*models.User)

	return nil
}

func (db *MockDB) DropSwipes(ctx context.Context) error {
	db.swipedUsers = make(map[uint64][]uint64)

	return nil
}

func (db *MockDB) CreateTag(ctx context.Context, tag_name string) error {
	newID := uint64(len(db.tags) + 1)
	db.tags[newID] = tag_name
	return nil
}

func (db *MockDB) GetTags(ctx context.Context) map[uint64]string {
	return db.tags
}
