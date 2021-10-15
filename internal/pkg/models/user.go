package models

import (
	"context"
	"crypto/md5"
	"dripapp/Models"
	"errors"
	"fmt"
)

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

func NewUser(id uint64, email string, password string) *User {
	return &User{ID: id, Email: email, Password: hashPassword(password)}
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid password")
)

// ArticleUsecase represent the article's usecases
type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	GetByTitle(ctx context.Context, title string) (User, error)
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

// ArticleRepository represent the article's repository contract
type UserRepository interface {
	GetUser(email string) (Models.User, error)
	GetUserByID(userID uint64) (Models.User, error)
	CreateUser(logUserData Models.LoginUser) (Models.User, error)
	AddSwipedUsers(currentUserId, swipedUserId uint64) error
	GetNextUserForSwipe(currentUserId uint64) (Models.User, error)
	UpdateUser(newUserData Models.User) error
	GetTags() map[uint64]string
}
