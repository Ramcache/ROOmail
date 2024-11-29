package router

import (
	"ROOmail/config"
	"ROOmail/internal/handlers/auth"
	"ROOmail/internal/handlers/handle"
	"ROOmail/internal/handlers/tasks"
	"ROOmail/internal/middleware"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(db *sql.DB, cfg config.Config) http.Handler {
	r := mux.NewRouter()

	// Маршруты для аутентификации
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")

	// Маршруты для задач
	taskService := tasks.NewTaskService(db)
	taskHandler := tasks.NewTaskHandler(taskService)

	// Защищённые маршруты для задач
	protectedRoutes := r.PathPrefix("/tasks").Subrouter()
	protectedRoutes.Use(middleware.JWTMiddleware)
	protectedRoutes.HandleFunc("", taskHandler.CreateTaskHandler).Methods("POST")
	protectedRoutes.HandleFunc("", taskHandler.GetTasksHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", taskHandler.GetTaskByIDHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", taskHandler.UpdateTaskHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", taskHandler.DeleteTaskHandler).Methods("DELETE")
	protectedRoutes.HandleFunc("/upload", taskHandler.UploadFileHandler).Methods("POST")
	protectedRoutes.HandleFunc("/download/{fileID}", taskHandler.DownloadFileHandler).Methods("GET")
	// Общедоступные маршруты
	usersService := handle.NewUsersService(db)
	usersHandler := handle.NewUsersHandler(usersService)

	r.HandleFunc("/users_list", usersHandler.UsersSelectHandler).Methods("GET")

	// Swagger-документация
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://chechenmail.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsHandler.Handler(r)
}
