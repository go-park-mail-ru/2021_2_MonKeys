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

type Image struct {
	img     image.Image
	gif     *gif.GIF
	imgType string
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

func (fm FileManager) SaveUserPhoto(user models.User, file io.Reader, fileName string) (path string, err error) {
	img, err := decodeImg(file, fileName)
	if err != nil {
		return
	}

	newPhoto, err := getNewFilename()
	if err != nil {
		return
	}

	path = fmt.Sprintf("%s/%s.webp", fm.getPathToUserPhoto(user), newPhoto)

	err = img.saveAsWebp(path)
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

func getNewFilename() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", u.String()), nil
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

func decodeImg(file io.Reader, fileName string) (img Image, err error) {
	fileType, err := getFileType(fileName)
	if err != nil {
		return
	}

	switch fileType {
	case "jpeg":
		img.img, err = jpeg.Decode(file)
	case "jpg":
		img.img, err = jpeg.Decode(file)
	case "png":
		img.img, err = png.Decode(file)
	case "gif":
		img.gif, err = gif.DecodeAll(file)
	case "webp":
		img.img, err = webp.Decode(file)
	default:
		err = errors.New("Unsupported file type")
	}

	img.imgType = fileType

	return
}

func (img Image) saveAsWebp(path string) error {
	fileOnDisk, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fileOnDisk.Close()

	if img.imgType == "gif" {
		err = gif.EncodeAll(fileOnDisk, img.gif)
	} else {
		err = webp.Encode(fileOnDisk, img.img, nil)
	}
	if err != nil {
		return err
	}

	return nil
}

func (img Image) getType() string {
	return img.imgType
}
