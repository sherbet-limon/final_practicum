package api

import (
	"net/http"
	// "go1f/pkg/db"
)

func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
}
