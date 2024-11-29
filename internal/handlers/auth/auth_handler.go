package auth

import (
	"ROOmail/pkg/utils"
	"encoding/json"
	"net/http"
	"strings"
)

var authService = NewAuthService()

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// LoginHandler выполняет вход пользователя
// @Summary Вход пользователя
// @Description Аутентификация пользователя и возвращение JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} string "Некорректный запрос"
// @Failure 401 {object} string "Неверное имя пользователя или пароль"
// @Failure 500 {object} string "Ошибка при генерации токена"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	user, err := authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// LogoutHandler выполняет выход пользователя
// @Summary Выход пользователя
// @Description Выход пользователя и отзыв JWT токена
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Success 303 "Перенаправление на страницу входа"
// @Failure 401 {object} string "Требуется заголовок авторизации"
// @Failure 401 {object} string "Некорректный формат заголовка авторизации"
// @Router /logout [get]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Требуется заголовок авторизации", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Некорректный формат заголовка авторизации", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	authService.RevokeToken(token)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
