package Models

import (
	"crypto/md5"
	"errors"
	"fmt"
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

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func MakeUser(id uint64, email, password string) User {
	return User{ID: id, Email: email, Password: hashPassword(password)}
}

func (user User) IsEmpty() bool {
	return len(user.Email) == 0
}

func (user User) IsCorrectPassword(password string) bool {
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

func (user *User) FillProfile(newUserData User) (err error) {
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
