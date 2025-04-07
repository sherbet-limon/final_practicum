package main

import (
	"log"
	"net/http"

)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))

	log.Println("Сервер запущен на http://localhost:7540")
	log.Fatal(http.ListenAndServe(":7540", nil))
}