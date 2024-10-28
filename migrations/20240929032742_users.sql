-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    status INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE contacts;
-- +goose StatementEnd
