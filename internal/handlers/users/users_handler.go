package users

import (
	"ROOmail/internal/models"
	_ "ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"ROOmail/pkg/utils/jwt_token"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	service *UserService
	log     logger.Logger
}

func NewUsersHandler(service *UserService, log logger.Logger) *UserHandler {
	return &UserHandler{service: service,
		log: log,
	}
}

// AddUserHandler обрабатывает запрос на добавление нового пользователя в базу данных.
// @Summary Добавить нового пользователя
// @Description Добавляет нового пользователя в базу данных с заданными именем, паролем и ролью.
// @Tags пользователи
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} map[string]interface{} "Сообщение об успешном добавлении и ID нового пользователя"
// @Failure 400 {object} string "Некорректные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на добавление нового пользователя")

	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Некорректный JSON при добавлении пользователя", err)
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		h.log.Error("Некорректные данные пользователя: отсутствуют обязательные поля")
		http.Error(w, "Некорректные данные: имя пользователя, пароль и роль обязательны", http.StatusBadRequest)
		return
	}

	userID, err := h.service.AddUser(r.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		h.log.Error("Не удалось добавить пользователя в базу данных", err)
		http.Error(w, "Не удалось добавить пользователя", http.StatusInternalServerError)
		return
	}

	h.log.Info("Пользователь успешно добавлен ", "userID: ", userID)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Пользователь успешно добавлен ", "user_id": %d}`, userID)))
}

// UsersSelectHandler
// @Summary      Получить список пользователей
// @Description  Возвращает список пользователей с возможностью фильтрации по имени пользователя.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username query string false "Фильтр по имени пользователя (поддерживает подстроку)"
// @Success      200 {array} models.UsersList
// @Failure      401 {object} map[string]string "Ошибка авторизации"
// @Failure      500 {object} map[string]string "Ошибка получения пользователей"
// @Router       /admin/users_list [get]
func (h *UserHandler) UsersSelectHandler(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value("user").(*jwt_token.Claims)
	if !ok {
		h.log.Error("Не удалось извлечь информацию о пользователе из контекста")
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	usernameFilter := r.URL.Query().Get("username")
	h.log.Info("Запрос списка пользователей. Фильтр по имени пользователя: ", usernameFilter)

	users, err := h.service.GetUsers(usernameFilter)
	if err != nil {
		h.log.Error("Ошибка получения пользователей: ", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Список пользователей успешно получен пользователем: %s", userClaims.Username)
	utils.RespondJSON(w, http.StatusOK, users)
}
