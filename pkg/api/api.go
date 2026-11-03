package api

import (
	"TODO_List/pkg/db"
	"encoding/json"
	"net/http"
	"time"
)

const formatTime = "20060102"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", nextDateHandler)
	mux.HandleFunc("POST /api/task", taskPostHandler)
	mux.HandleFunc("GET /api/tasks", tasksHandler)
	mux.HandleFunc("GET /api/task", taskGetHandler)
	mux.HandleFunc("PUT /api/task", taskPutHandler)
	mux.HandleFunc("POST /api/task/done", taskDoneHandler)
	mux.HandleFunc("DELETE /api/task", taskDeleteHandler)
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format(formatTime)
	}

	now, err := time.Parse(formatTime, nowStr)
	if err != nil {
		http.Error(w, "invalid now date", http.StatusBadRequest)
		return
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}

func taskPostHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if t.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJson(w, map[string]string{"error": "db error"}, http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]any{"id": id}, http.StatusOK)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	list, err := db.Tasks(50)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if list == nil {
		list = []*db.Task{}
	}
	writeJson(w, map[string][]*db.Task{"tasks": list}, http.StatusOK)
}

func taskGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"}, http.StatusBadRequest)
		return
	}
	writeJson(w, t, http.StatusOK)
}

func taskPutHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if t.ID == "" || t.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок или идентификатор"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&t); err != nil {
		writeJson(w, map[string]string{"error": "update error"}, http.StatusBadRequest)
		return
	}
	writeJson(w, struct{}{}, http.StatusOK)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "task not found"}, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, struct{}{}, http.StatusOK)
		return
	}

	nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err = db.UpdateDate(nextDate, task.ID); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, struct{}{}, http.StatusOK)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, struct{}{}, http.StatusOK)
}

func checkDate(t *db.Task) error {
	today := time.Now().Truncate(24 * time.Hour)

	if t.Date == "" {
		t.Date = today.Format(formatTime)
	}

	taskDate, err := time.Parse(formatTime, t.Date)
	if err != nil {
		return err
	}

	if taskDate.Before(today) {
		if t.Repeat == "" {
			t.Date = today.Format(formatTime)
		} else {
			next, err := NextDate(today, t.Date, t.Repeat)
			if err != nil {
				return err
			}
			t.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
