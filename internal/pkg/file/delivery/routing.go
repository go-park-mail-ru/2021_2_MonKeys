package delivery

import (
	"dripapp/internal/pkg/file"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetFileRouting(router *mux.Router, fm file.FileManager) {
	router.PathPrefix(fmt.Sprintf("/%s/", fm.RootFolder)).Handler(
		http.StripPrefix(fmt.Sprintf("/%s/", fm.RootFolder),
			http.FileServer(http.Dir("./" + fm.RootFolder))))

	http.Handle("/", router)
}
