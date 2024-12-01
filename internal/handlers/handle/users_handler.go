package handle

import (
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"net/http"
)

type UsersHandler struct {
	service *UsersService
	log     logger.Logger
}

func NewUsersHandler(service *UsersService, log logger.Logger) *UsersHandler {
	return &UsersHandler{
		service: service,
		log:     log,
	}

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
	h.log.Info("Запрос списка пользователей. Фильтр по имени пользователя: ", usernameFilter)

	users, err := h.service.GetUsers(usernameFilter)
	if err != nil {
		h.log.Error("Ошибка получения пользователей: ", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	h.log.Info("Список пользователей успешно получен")
	utils.RespondJSON(w, http.StatusOK, users)
}
