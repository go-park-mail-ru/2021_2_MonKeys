package file

import (
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/hasher"
	"errors"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
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

	if ok, err := fm.isNotExists(fm.RootFolder); ok {
		if err != nil {
			return nil, err
		}

		err = fm.createFolder(fm.RootFolder)
		if err != nil {
			return nil, err
		}
	}

	if ok, err := fm.isNotExists(fm.PhotoFolder); ok {
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

//func (fm FileManager) SaveUserPhoto(user models.User, file io.Reader, fileName string) (path string, err error) {
//	fileType, err := getFileType(fileName)
//	if err != nil {
//		return
//	}
//
//	err = validImgType(fileType)
//	if err != nil {
//		return
//	}
//
//	newPhoto, err := getNewFilename(fileType)
//	if err != nil {
//		return
//	}
//
//	path = fmt.Sprintf("%s/%s", fm.getPathToUserPhoto(user), newPhoto)
//
//	err = fm.save(file, path)
//	return
//}
func (fm FileManager) SaveUserPhoto(user models.User, file io.Reader, fileName string) (path string, err error) {
	fileType, err := getFileType(fileName)
	if err != nil {
		return
	}

	img, err := fm.decodeImg(file, fileType)
	if err != nil {
		return
	}

	newPhoto, err := getNewFilename("webp")
	if err != nil {
		return
	}

	path = fmt.Sprintf("%s/%s", fm.getPathToUserPhoto(user), newPhoto)

	err = fm.saveAsWebp(img, path)
	return
}

func (FileManager) Delete(filePath string) error {
	return os.Remove(filePath)
}

func (FileManager) createFolder(title string) error {
	return os.Mkdir(title, os.ModePerm)
}

func (FileManager) save(file io.Reader, path string) error {
	fileOnDisk, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fileOnDisk.Close()

	_, err = io.Copy(fileOnDisk, file)
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

	fileType := separatedFilename[len(separatedFilename)-1]

	return strings.ToLower(fileType), nil
}

func getNewFilename(fileType string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", u.String(), fileType), nil
}

func (FileManager) isNotExists(path string) (bool, error) {
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
	if fileType != "png" &&
		fileType != "jpg" &&
		fileType != "jpeg" &&
		fileType != "gif" &&
		fileType != "webp" {
		return errors.New("wrong file type")
	}

	return nil
}

func (FileManager) decodeImg(file io.Reader, fileType string) (img image.Image, err error) {
	switch fileType {
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "jpg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	case "webp":
		img, err = webp.Decode(file)
	default:
		err = errors.New("Unsupported file type")
	}
	return
}

func (FileManager) saveAsWebp(img image.Image, path string) error {
	fileOnDisk, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fileOnDisk.Close()

	err = webp.Encode(fileOnDisk, img, nil)
	if err != nil {
		return err
	}

	return nil
}
