-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL 
)
-- +goose StatementEnd

INSERT INTO users (username, email) VALUES ('User 1', 'username1@gmail.com');
INSERT INTO users (username, email) VALUES ('User 2', 'username2@gmail.com');
INSERT INTO users (username, email) VALUES ('User 3', 'username3@gmail.com');

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
