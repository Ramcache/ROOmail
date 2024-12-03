package auth

import (
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils"
	"ROOmail/pkg/utils/jwt"
	"encoding/json"
	"net/http"
	"strings"
)

var (
	authService = AuthServiceInstance()
	log         = logger.NewZapLogger()
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Используем контекст из запроса
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Ошибка декодирования тела запроса: ", err)
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректный запрос"})
		return
	}

	log.Info("Попытка входа пользователя: ", req.Username)
	user, err := authService.AuthenticateUser(ctx, req.Username, req.Password)
	if err != nil {
		log.Warn("Неудачная попытка входа пользователя: ", req.Username)
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверное имя пользователя или пароль"})
		return
	}

	token, err := jwt.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		log.Error("Ошибка генерации токена для пользователя: ", req.Username, " - ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при генерации токена"})
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Используем контекст из запроса
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn("Попытка выхода без заголовка авторизации")
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Требуется заголовок авторизации"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Warn("Некорректный формат заголовка авторизации: ", authHeader)
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Некорректный формат заголовка авторизации"})
		return
	}

	token := parts[1]
	if err := authService.RevokeToken(ctx, token); err != nil {
		log.Error("Ошибка отзыва токена: ", token, " - ", err)
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при отзыве токена"})
		return
	}

	log.Info("Пользователь вышел, токен отозван: ", token)
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Выход выполнен успешно"})
}
