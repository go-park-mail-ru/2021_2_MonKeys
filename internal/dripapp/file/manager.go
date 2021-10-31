package file

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"fmt"
	"io"
	"os"
	"strconv"
)

type FileManager struct {
	RootFolder  string
	PhotoFolder string
}

func NewFileManager(config configs.FileStorageConfig) (fm *FileManager, err error) {
	rootFolder := config.RootFolder
	photoFolder := fmt.Sprintf("%s/%s", rootFolder, config.ProfilePhotoPath)

	fm = &FileManager{
		RootFolder:  rootFolder,
		PhotoFolder: photoFolder,
	}

	err = fm.createFolder(fm.RootFolder)

	err = fm.createFolder(fm.PhotoFolder)

	return fm, err
}

func (fm FileManager) CreateFoldersForNewUser(user models.User) error {
	return fm.createFolder(fm.getPathToUserPhoto(user))
}

func (fm FileManager) SaveUserPhoto(user models.User, file io.Reader) (path string, err error) {
	path, err = fm.getPathToNewPhoto(user)
	if err != nil {
		return
	}

	err = fm.save(path, file)
	return
}

func (FileManager) Delete(filePath string) error {
	return os.Remove(filePath)
}

func (FileManager) createFolder(title string) error {
	return os.Mkdir(title, 0777)
}

func (FileManager) save(path string, file io.Reader) error {
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

func (fm FileManager) getPathToUserPhoto(user models.User) string {
	return fmt.Sprintf("%s/%s", fm.PhotoFolder, user.Email)
}

func (fm FileManager) getPathToNewPhoto(user models.User) (pathToNewPhoto string, err error) {
	newPhoto, err := fm.createNameToNewPhoto(user)
	if err != nil {
		return
	}

	pathToNewPhoto = fmt.Sprintf("%s/%s", fm.getPathToUserPhoto(user), newPhoto)

	return
}

func (fm FileManager) createNameToNewPhoto(user models.User) (string, error) {
	lastPhoto := user.GetLastPhoto()
	if lastPhoto == "" {
		return "1.png", nil
	}

	numStr := lastPhoto[len(fm.getPathToUserPhoto(user))+1 : len(lastPhoto)-4]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(num+1) + ".png", nil
}
