-- +goose Up
CREATE TABLE task_assignments (
                                  id SERIAL PRIMARY KEY,
                                  task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
                                  user_id INTEGER NOT NULL,
                                  assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE task_assignments;
