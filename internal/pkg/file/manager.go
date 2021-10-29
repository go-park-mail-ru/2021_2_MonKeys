package file

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/pkg/models"
	"fmt"
	"io"
	"os"
)

type FileManager struct {
	RootFolder  string
	PhotoFolder string
}

func NewFileManager(config configs.FileStorageConfig) (*FileManager, error) {
	fm := FileManager{
		RootFolder:  config.RootFolder,
		PhotoFolder: config.RootFolder + "/" + config.ProfilePhotoPath,
	}

	err := fm.CreateFolder(fm.RootFolder)

	return &fm, err
}

func (FileManager) CreateFolder(title string) error {
	return os.Mkdir(title, 0777)
}

func (fm FileManager) CreateFoldersForNewUser(user models.User) error {
	return os.Mkdir(fmt.Sprintf("%s/%s", fm.PhotoFolder, user.Email), 0777)
}

func (FileManager) Save(ctx context.Context, path string, file io.Reader) error {
	saved, err := os.Create(path)
	if err != nil {
		return err
	}
	defer saved.Close()

	_, err = io.Copy(saved, file)
	if err != nil {
		return err
	}

	return nil
}

func (fm FileManager) getNewPathForUser(user models.User) string {
	return fmt.Sprintf("%s/%s/%s", fm.PhotoFolder, user.Email, user.GetNameToNewPhoto())
}

func (fm FileManager) SaveUserPhoto(ctx context.Context, user models.User, file io.Reader) (path string, err error) {
	path = fm.getNewPathForUser(user)
	err = fm.Save(ctx, path, file)
	return
}

func (FileManager) Delete(ctx context.Context, filePath string) error {
	return os.Remove(filePath)
}
