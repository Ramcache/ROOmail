package router

import (
	"ROOmail/config"
	"ROOmail/internal/handlers"
	"ROOmail/internal/handlers/auth"
	"ROOmail/internal/handlers/file"
	"ROOmail/internal/handlers/tasks"
	"ROOmail/internal/handlers/users"
	"ROOmail/pkg/logger"
	"ROOmail/pkg/utils/jwt_token"
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

	// Регистрация маршрутов работы с файлами
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

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(jwt_token.JWTMiddleware)
	adminRouter.Use(jwt_token.RoleMiddleware("admin"))
	adminRouter.HandleFunc("/logs/list", handlers.ListLogsHandler).Methods("GET")
	adminRouter.HandleFunc("/logs/{filename}", handlers.LogsHandler).Methods("GET")
}

// Регистрация маршрутов для задач
func registerTaskRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	taskService := tasks.NewTaskService(db)
	taskHandler := tasks.NewTaskHandler(taskService, log)

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(jwt_token.JWTMiddleware)
	adminRouter.Use(jwt_token.RoleMiddleware("admin"))
	adminRouter.HandleFunc("/tasks/create", taskHandler.CreateTaskHandler).Methods("POST") //1
	adminRouter.HandleFunc("/tasks/update/{id}", taskHandler.UpdateTaskHandler).Methods("PUT")
	adminRouter.HandleFunc("/tasks/update/{id}", taskHandler.PatchTaskHandler).Methods("PATCH")
	adminRouter.HandleFunc("/tasks/delete/{id}", taskHandler.DeleteTaskHandler).Methods("DELETE")

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(jwt_token.JWTMiddleware)
	userRouter.Use(jwt_token.RoleMiddleware("users"))
	userRouter.HandleFunc("/tasks/get/{id}", taskHandler.GetTasksHandler).Methods("GET")
}

// Регистрация маршрутов для пользователей
func registerUserRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	usersService := users.NewUsersService(db)
	usersHandler := users.NewUsersHandler(usersService, log)

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(jwt_token.JWTMiddleware)
	adminRouter.Use(jwt_token.RoleMiddleware("admin"))
	adminRouter.HandleFunc("/users_list", usersHandler.UsersSelectHandler).Methods("GET")
	adminRouter.HandleFunc("/users/add", usersHandler.AddUserHandler).Methods("POST")
	adminRouter.HandleFunc("/users/delete/{id}", usersHandler.DeleteUserHandler).Methods("DELETE")

}

func registerFIleRoutes(r *mux.Router, db *pgxpool.Pool, log logger.Logger) {
	fileService := file.NewFileService("./uploads", db)
	fileHandler := file.NewFileHandler(fileService, log)

	fileRouter := r.PathPrefix("/admin").Subrouter()
	fileRouter.Use(jwt_token.JWTMiddleware)
	fileRouter.Use(jwt_token.RoleMiddleware("admin"))
	fileRouter.HandleFunc("/file/upload", fileHandler.UploadFileHandler).Methods("POST")
}
