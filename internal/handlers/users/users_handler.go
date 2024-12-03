package users

import (
	"ROOmail/internal/models"
	_ "ROOmail/internal/models"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"ROOmail/pkg/utils/jwt_token"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные пользователя"
// @Success 201 {object} map[string]interface{} "Сообщение об успешном добавлении и ID нового пользователя"
// @Failure 400 {object} string "Некорректные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /admin/users/add [post]
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"Пользователь успешно добавлен ", "user_id": %d}`, userID)))
}

// DeleteUserHandler обрабатывает запрос на удаление пользователя по его ID.
// @Summary Удалить пользователя
// @Description Удаляет пользователя из базы данных по его идентификатору (ID).
// @Tags users
// @Param id path int true "ID пользователя"
// @Success 204 "Пользователь успешно удалён"
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 404 {object} string "Пользователь не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /admin/users/delete/{id} [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Получен запрос на удаление пользователя")

	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.log.Error("Некорректный идентификатор пользователя", err)
		http.Error(w, "Некорректный запрос: некорректный идентификатор пользователя", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), userID)
	if err != nil {
		h.log.Error("Не удалось удалить пользователя", err)
		http.Error(w, fmt.Sprintf("Не удалось удалить пользователя с ID %d", userID), http.StatusInternalServerError)
		return
	}

	h.log.Info("Пользователь успешно удалён ", "userID: ", userID)

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf(`{"Пользователь успешно удален ", "user_id": %d}`, userID)))
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
