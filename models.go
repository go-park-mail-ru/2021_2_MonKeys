package main

import (
	"crypto/md5"
	"fmt"
)

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type User struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"-"`
	Age         uint     `json:"age"`
	Description string   `json:"description"`
	ImgSrc      string   `json:"imgSrc"`
	Tags        []string `json:"tags"`
}

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func makeUser(id uint64, email, password string) User {
	return User{ID: id, Email: email, Password: hashPassword(password)}
}

func (user User) isEmpty() bool {
	return len(user.Email) == 0
}

func (user User) isCorrectPassword(password string) bool {
	return user.Password == hashPassword(password)
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SwipedUser struct {
	Id uint64 `json:"id"`
}
