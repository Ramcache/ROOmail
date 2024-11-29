-- +goose Up
CREATE TABLE IF NOT EXISTS tasks
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER      NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    due_date    DATE,
    file        VARCHAR(255),
    priority    INTEGER      NOT NULL,
    schools     VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS tasks;
