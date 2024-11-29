package utils

import (
	"ROOmail/pkg/logger"
	"encoding/json"
	"net/http"
)

var log = logger.GetLogger()

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Error("Ошибка маршализации JSON ответа: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
