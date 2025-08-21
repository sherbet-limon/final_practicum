package api

import (
	"net/http"
)
// обработчик методов post, put etc
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodGet:	
		GetTaskHandler(w, r)
	case http.MethodDelete:
		TaskDelHandler(w, r)
	}
}
