package router

import (
	"ROOmail/config"
	"ROOmail/internal/handlers/auth"
	"ROOmail/internal/handlers/file"
	"ROOmail/internal/handlers/tasks"
	"ROOmail/internal/handlers/users"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils/JWT"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func InitRouter(db *pgxpool.Pool, cfg config.Config) http.Handler {
	r := mux.NewRouter()
	log := logger.NewZapLogger()

	// Регистрация маршрутов аутентификации
	registerAuthRoutes(r, log)

	// Регистрация маршрутов задач
	registerTaskRoutes(r, db, log)

	// Регистрация маршрутов пользователей
	registerUserRoutes(r, db, log)

	registerFIleRoutes(r, db, log)

	// Swagger-документация
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// CORS настройки
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://chechenmail.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(r)
}

// Регистрация маршрутов для аутентификации
func registerAuthRoutes(r *mux.Router, log logger.Logger) {
	r.HandleFunc("/auth/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/logout", auth.LogoutHandler).Methods("POST")
}

// Регистрация маршрутов для задач
func registerTaskRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	taskService := tasks.NewTaskService(db)
	taskHandler := tasks.NewTaskHandler(taskService, log)

	protectedRouter := r.PathPrefix("/admin").Subrouter()
	protectedRouter.Use(JWT.JWTMiddleware)
	protectedRouter.Use(JWT.RoleMiddleware("admin"))
	protectedRouter.HandleFunc("/tasks/create", taskHandler.CreateTaskHandler).Methods("POST") //1
}

// Регистрация маршрутов для пользователей
func registerUserRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	usersService := users.NewUsersService(db)
	usersHandler := users.NewUsersHandler(usersService, log)

	userRouter := r.PathPrefix("/admin").Subrouter()
	userRouter.Use(JWT.JWTMiddleware)
	userRouter.Use(JWT.RoleMiddleware("admin"))
	userRouter.HandleFunc("/users_list", usersHandler.UsersSelectHandler).Methods("GET")
}

func registerFIleRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	fileService := file.NewFileService("./uploads", db)
	fileHandler := file.NewFileHandler(fileService, log)

	fileRouter := r.PathPrefix("/admin").Subrouter()
	fileRouter.Use(JWT.JWTMiddleware)
	fileRouter.Use(JWT.RoleMiddleware("admin"))
	fileRouter.HandleFunc("/file/upload", fileHandler.UploadFileHandler).Methods("POST")
}
