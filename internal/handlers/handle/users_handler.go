package handle

import (
	"encoding/json"
	"log"
	"net/http"
)

type UsersHandler struct {
	service *UsersService
}

func NewUsersHandler(service *UsersService) *UsersHandler {
	return &UsersHandler{service: service}
}
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// UsersSelectHandler
// @Summary      Получить список пользователей
// @Description  Возвращает список пользователей с возможностью фильтрации по имени пользователя.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        username query string false "Фильтр по имени пользователя (поддерживает подстроку)"
// @Success      200 {array} models.UsersList
// @Failure      500 {object} map[string]string
// @Router       /users_list [get]
func (h *UsersHandler) UsersSelectHandler(w http.ResponseWriter, r *http.Request) {
	usernameFilter := r.URL.Query().Get("username")

	users, err := h.service.GetUsers(usernameFilter)
	if err != nil {
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		log.Println("Ошибка бизнес-логики:", err)
		return
	}

	RespondJSON(w, http.StatusOK, users)
}
