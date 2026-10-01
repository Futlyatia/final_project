package api

import (
	"net/http"
	"time"

	"final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var date string

	if search != "" {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			date = parsedDate.Format("20060102")
			search = ""
		}
	}

	tasks, err := db.Tasks(50, search, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
