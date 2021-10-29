package models

import (
	"context"
	"io"
)

type FileRepository interface {
	CreateFolder(title string) error
	CreateFoldersForNewUser(user User) error
	Save(ctx context.Context, path string, file io.Reader) error
	SaveUserPhoto(ctx context.Context, user User, file io.Reader) (path string, err error)
	Delete(ctx context.Context, filePath string) error
}
