package delivery

import (
	"dripapp/internal/dripapp/file"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SetFileRouting(router *mux.Router, fm file.FileManager) {
	router.PathPrefix(fmt.Sprintf("/%s/", fm.RootFolder)).Handler(
		http.StripPrefix(fmt.Sprintf("/%s/", fm.RootFolder),
			http.FileServer(http.Dir("./"+fm.RootFolder))))

	http.Handle("/", router)
}
