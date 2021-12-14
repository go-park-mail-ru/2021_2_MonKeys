package models

import (
	"io"
)

//easyjson:json
type Photo struct {
	Path string `json:"photo"`
}

func (user User) GetLastPhoto() string {
	if len(user.Imgs) == 0 {
		return ""
	}

	return user.Imgs[len(user.Imgs)-1]
}

func (user *User) AddNewPhoto(photoPath string) {
	user.Imgs = append(user.Imgs, photoPath)
}

func (user *User) DeletePhoto(photo Photo) (err error) {
	var photos []string

	err = ErrNoSuchPhoto
	for _, currPhoto := range user.Imgs {
		if currPhoto != photo.Path {
			photos = append(photos, currPhoto)
			err = nil
		}
	}

	user.Imgs = photos

	return
}

type FileRepository interface {
	CreateFoldersForNewUser(user User) error
	SaveUserPhoto(user User, file io.Reader, fileName string) (path string, err error)
	Delete(filePath string) error
}
