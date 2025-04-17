package db

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    // обработка других методов будет добавлена на следующих шагах
    case http.MethodPost:
        func addTaskHandler(w, r)
    }
} 
func AddTask(task *Task) (int64, error){

}