package server

import (
	"fmt"
	"go1f/pkg/api"
	"log"
	"net/http"
	"os"
)

// запускает локальный сервер на порту 7540
func Run() error {
	api.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "7540" 
	}
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Println("Сервер запущен на http://localhost:7540")
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
