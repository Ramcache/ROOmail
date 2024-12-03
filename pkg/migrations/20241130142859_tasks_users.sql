-- +goose Up
CREATE TABLE tasks_users (
                             task_id INT NOT NULL,
                             user_id INT NOT NULL,
                             assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             PRIMARY KEY (task_id, user_id),
                             FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE IF EXISTS task_recipients;
