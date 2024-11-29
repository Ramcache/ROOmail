package models

type Task struct {
	ID          string   `json:"id"`
	UserID      int      `json:"user_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	DueDate     string   `json:"due_date"`
	File        []string `json:"file"`
	Priority    string   `json:"priority"`
	Schools     string   `json:"schools"`
}
