package models

import (
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
	ImgSrc      string   `json:"imgSrc,omitempty"`
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
	Id      uint64 `json:"id"`
	TagText string `json:"tagText"`
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
	user.ImgSrc = newUserData.ImgSrc
	user.Tags = newUserData.Tags

	return nil
}

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
)

// ArticleUsecase represent the article's usecases
type UserUsecase interface {
	CurrentUser(w http.ResponseWriter, r *http.Request)
	EditProfileHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	SignupHandler(w http.ResponseWriter, r *http.Request)
	NextUserHandler(w http.ResponseWriter, r *http.Request)
	GetAllTags(w http.ResponseWriter, r *http.Request)
}

// ArticleRepository represent the article's repository contract
type UserRepository interface {
	GetUser(email string) (*User, error)
	GetUserByID(userID uint64) (*User, error)
	CreateUser(logUserData *LoginUser) (*User, error)
	UpdateUser(newUserData *User) error
	AddSwipedUsers(currentUserId, swipedUserId uint64) error
	GetNextUserForSwipe(currentUserId uint64) (*User, error)
	IsSwiped(userID, swipedUserID uint64) bool
	CreateUserAndProfile(user *User)
	DropUsers()
	DropSwipes()
	CreateTag(text string)
	GetTags() map[uint64]string
}
