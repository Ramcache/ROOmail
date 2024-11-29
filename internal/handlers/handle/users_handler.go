package handle

import (
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"net/http"
)

var log = logger.GetLogger() // Получаем глобальный логгер

type UsersHandler struct {
	service *UsersService
}

func NewUsersHandler(service *UsersService) *UsersHandler {
	return &UsersHandler{service: service}
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
	log.Info("Запрос списка пользователей. Фильтр по имени пользователя: ", usernameFilter)

	users, err := h.service.GetUsers(usernameFilter)
	if err != nil {
		log.Error("Ошибка получения пользователей: ", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	log.Info("Список пользователей успешно получен")
	utils.RespondJSON(w, http.StatusOK, users)
}
