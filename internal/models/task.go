package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	UserIDs     []int  `json:"user_ids"`
	FilePath    string `json:"file_path,omitempty"`
	CreatedBy   int    `json:"created_by"`
}
