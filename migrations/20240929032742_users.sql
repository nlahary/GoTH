-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE
)
-- +goose StatementEnd

INSERT INTO contacts (username, email) VALUES ('User 1', 'username1@gmail.com');
INSERT INTO contacts (username, email) VALUES ('User 2', 'username2@gmail.com');
INSERT INTO contacts (username, email) VALUES ('User 3', 'username3@gmail.com');

-- +goose Down
-- +goose StatementBegin
DROP TABLE contacts;
-- +goose StatementEnd
