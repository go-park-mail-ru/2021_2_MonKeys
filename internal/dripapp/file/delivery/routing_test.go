package delivery

import (
	"dripapp/internal/dripapp/file"
	"github.com/gorilla/mux"
	"testing"
)

func TestSetFileRouting(t *testing.T) {
	fm := file.FileManager{
		RootFolder:  "",
		PhotoFolder: "",
	}
	SetFileRouting(mux.NewRouter(), fm)
}
