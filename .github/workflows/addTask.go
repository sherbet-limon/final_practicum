// package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"go1f/pkg/db"
// 	"net/http"
// )

// func addTaskHandler(w http.ResponseWriter, r *http.Request) {
// 	var task db.Task
// 	var buf bytes.Buffer
// 	_, err := buf.ReadFrom(r.Body)
// 	if err != nil {
// 		http.Error(w, "ошибка чтения", http.StatusBadRequest)
// 		return
// 	}
// 	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
// 		http.Error(w, "ошибка десериализации", http.StatusBadRequest)
// 		return
// 	}
//     if task.Title == "" {
//         http.Error(w, "заголовок не должен быть пустым", http.StatusBadRequest)
// }
// date, err := CheckDate(task.Date)
// }
