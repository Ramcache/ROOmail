package auth

import (
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"encoding/json"
	"net/http"
	"strings"
)

var (
	authService = NewAuthService()
	log         = logger.GetLogger()
)

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
		log.Error("Ошибка декодирования тела запроса: ", err)
		utils.RespondJSON(w, http.StatusBadRequest, "Некорректный запрос")
		return
	}

	log.Info("Попытка входа пользователя: ", req.Username)
	user, err := authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log.Warn("Неудачная попытка входа пользователя: ", req.Username)
		utils.RespondJSON(w, http.StatusUnauthorized, "Неверное имя пользователя или пароль")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Error("Ошибка генерации токена для пользователя: ", req.Username, " - ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "Ошибка при генерации токена")
		return
	}

	resp := LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}
	log.Info("Успешный вход пользователя: ", user.Username)
	utils.RespondJSON(w, http.StatusOK, resp)
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
		log.Warn("Попытка выхода без заголовка авторизации")
		utils.RespondJSON(w, http.StatusUnauthorized, "Требуется заголовок авторизации")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Warn("Некорректный формат заголовка авторизации: ", authHeader)
		utils.RespondJSON(w, http.StatusUnauthorized, "Некорректный формат заголовка авторизации")
		return
	}

	token := parts[1]
	authService.RevokeToken(token)
	log.Info("Пользователь вышел, токен отозван: ", token)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
