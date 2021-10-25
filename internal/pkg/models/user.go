package models

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type User struct {
	ID          uint64   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Password    string   `json:"-"`
	Date        string   `json:"date,omitempty"`
	Age         uint     `json:"age,omitempty"`
	Description string   `json:"description,omitempty"`
	Imgs        []string `json:"imgSrc,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SwipedUser struct {
	Id uint64 `json:"id"`
}

type Tag struct {
	Id       uint64 `json:"id"`
	Tag_Name string `json:"tagText"`
}

type Tags struct {
	AllTags map[uint64]Tag `json:"allTags"`
	Count   uint64         `json:"tagsCount"`
}

func NewUser(id uint64, email string, password string) *User {
	return &User{ID: id, Email: email, Password: hashPassword(password)}
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func (user *User) IsEmpty() bool {
	return len(user.Email) == 0
}

func (user *User) IsCorrectPassword(password string) bool {
	return user.Password == hashPassword(password)
}

func getAgeFromDate(date string) (uint, error) {
	birthday, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, errors.New("failed on userYear")
	}

	age := uint(time.Now().Year() - birthday.Year())
	if time.Now().YearDay() < birthday.YearDay() {
		age -= 1
	}

	return age, nil
}

func (user *User) FillProfile(newUserData *User) (err error) {
	user.Name = newUserData.Name
	user.Date = newUserData.Date
	user.Age, err = getAgeFromDate(newUserData.Date)
	if err != nil {
		return errors.New("failed to save age")
	}
	user.Date = newUserData.Date
	user.Description = newUserData.Description
	user.Imgs = newUserData.Imgs
	user.Tags = newUserData.Tags

	return nil
}

var (
	ErrNoUser             = errors.New("no user found")
	ErrBadPass            = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// ArticleUsecase represent the article's usecases
type UserUsecase interface {
	CurrentUser(c context.Context, r *http.Request) (User, int)
	EditProfile(c context.Context, newUserData User, r *http.Request) (User, int)
	Login(c context.Context, logUserData LoginUser, w http.ResponseWriter) (User, int)
	Logout(c context.Context, w http.ResponseWriter, r *http.Request) int
	Signup(c context.Context, logUserData LoginUser, w http.ResponseWriter) int
	NextUser(c context.Context, swipedUserData SwipedUser, r *http.Request) (User, int)
	GetAllTags(c context.Context, r *http.Request) (Tags, int)
}

// ArticleRepository represent the article's repository contract
type UserRepository interface {
	GetUser(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userID uint64) (*User, error)
	CreateUser(ctx context.Context, logUserData *LoginUser) (*User, error)
	UpdateUser(ctx context.Context, newUserData *User) error
	AddSwipedUsers(ctx context.Context, currentUserId uint64, swipedUserId uint64, type_name string) error
	GetNextUserForSwipe(ctx context.Context, currentUserId uint64) (User, error)
	IsSwiped(ctx context.Context, userID, swipedUserID uint64) (bool, error)
	CreateUserAndProfile(ctx context.Context, user *User) (User, error)
	DropUsers(ctx context.Context) error
	DropSwipes(ctx context.Context) error
	CreateTag(ctx context.Context, tag_name string) error
	GetTags(ctx context.Context) map[uint64]string
}
