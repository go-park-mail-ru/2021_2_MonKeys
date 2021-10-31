package file

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/hasher"
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"io"
	"os"
	"strings"
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

	if ok, err := isNotExists(fm.RootFolder); ok {
		if err != nil {
			return nil, err
		}

		err = fm.createFolder(fm.RootFolder)
		if err != nil {
			return nil, err
		}
	}

	if ok, err := isNotExists(fm.PhotoFolder); ok {
		if err != nil {
			return nil, err
		}

		err = fm.createFolder(fm.PhotoFolder)
		if err != nil {
			return nil, err
		}
	}

	return fm, err
}

func (fm FileManager) CreateFoldersForNewUser(user models.User) error {
	return fm.createFolder(fm.getPathToUserPhoto(user))
}

func (fm FileManager) SaveUserPhoto(user models.User, file io.Reader, fileName string) (path string, err error) {
	fileType, err := getFileType(fileName)
	if err != nil {
		return
	}

	err = validImgType(fileType)

	newPhoto, err := createNameToNewPhoto(fileType)
	if err != nil {
		return
	}

	path = fmt.Sprintf("%s/%s", fm.getPathToUserPhoto(user), newPhoto)

	err = fm.save(path, file)
	return
}

func (FileManager) Delete(filePath string) error {
	return os.Remove(filePath)
}

func (FileManager) createFolder(title string) error {
	return os.Mkdir(title, os.ModePerm)
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
	return fmt.Sprintf("%s/%s", fm.PhotoFolder, hasher.GetSha1([]byte(user.Email)))
}

func getFileType(fileName string) (string, error) {
	separatedFilename := strings.Split(fileName, ".")
	if len(separatedFilename) <= 1 {
		err := errors.New("bad filename")
		return "", err
	}

	return separatedFilename[len(separatedFilename)-1], nil
}

func createNameToNewPhoto(fileType string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", u.String(), fileType), nil
}

func isNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return false, nil
	}

	if os.IsNotExist(err) {
		return true, nil

	}

	return false, err
}

func validImgType(fileType string) error {
	fileTypeL := strings.ToLower(fileType)
	if fileTypeL != "png" &&
		fileTypeL != "jpg" &&
		fileTypeL != "jpeg" &&
		fileTypeL != "gif" {
		return errors.New("wrong file type")
	}

	return nil
}
