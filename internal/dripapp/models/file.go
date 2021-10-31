package models

import (
	"io"
)

type FileRepository interface {
	CreateFoldersForNewUser(user User) error
	SaveUserPhoto(user User, file io.Reader) (path string, err error)
	Delete(filePath string) error
}
