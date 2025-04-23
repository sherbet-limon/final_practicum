package api

import (
	"net/http"
	)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    // обработка других методов будет добавлена на следующих шагах
    case http.MethodPost:
		AddTaskHandler(w, r)
    }
} 
