-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status INTEGER NOT NULL DEFAULT 0
)
-- +goose StatementEnd

INSERT INTO contacts (username, email, status) VALUES ('User 1', 'username1@gmail.com', 1);
INSERT INTO contacts (username, email, status) VALUES ('User 2', 'username2@gmail.com', 1);
INSERT INTO contacts (username, email, status) VALUES ('GuestUser 3', 'username3@gmail.com');

-- +goose Down
-- +goose StatementBegin
DROP TABLE contacts;
-- +goose StatementEnd
