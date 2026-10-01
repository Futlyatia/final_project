package api

import (
	"net/http"
	"time"

	"final_project/pkg/db"
	"final_project/pkg/service"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан id задачи",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Если задача одноразовая,
	// после выполнения удаляем её.
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	// Получаем сегодняшнюю дату без времени.
	now := time.Now()
	now = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	// Для периодической задачи вычисляем
	// следующую дату выполнения.
	next, err := service.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	task.Date = next

	err = db.UpdateTask(task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
