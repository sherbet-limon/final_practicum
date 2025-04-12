package server

import (
	"fmt"
	"log"
	"net/http"

)
//локальный сервер на порту 7540
func Run() error {
	port := 7540
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Println("Сервер запущен на http://localhost:7540")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

