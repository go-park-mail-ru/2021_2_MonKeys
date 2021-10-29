package main

//import (
//	"log"
//	"net/http"
//)
//
//type MediaServerConfig struct {
//	Host string
//	Port string
//}
//
//var Server MediaServerConfig
//
//func init() {
//	Server = MediaServerConfig{
//		Host:     "127.0.0.1",
//		Port:     ":9999",
//	}
//}
//
//func main() {
//	staticHandler := http.StripPrefix(
//		"/media/",
//		http.FileServer(http.Dir("./media")),
//	)
//	http.Handle("/media/", staticHandler)
//
//	log.Println("starting server at ", Server.Port)
//
//	err := http.ListenAndServe(Server.Port, nil)
//	if err != nil {
//		return
//	}
//}
