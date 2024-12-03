package models

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	UserIDs     []int  `json:"user_ids"`
	FilePath    string `json:"file_path,omitempty"` // Путь к файлу, если есть
}
