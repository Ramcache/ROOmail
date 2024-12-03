-- +goose Up
CREATE TABLE IF NOT EXISTS task_recipients (
                                               task_id INTEGER NOT NULL,
                                               user_id INTEGER NOT NULL,
                                               PRIMARY KEY (task_id, user_id),
                                               FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE,
                                               FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS task_recipients;
