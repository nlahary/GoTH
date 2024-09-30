-- +goose Up
-- +goose StatementBegin
CREATE TABLE cartStatus (
    id INTEGER PRIMARY KEY,
    status VARCHAR(10) NOT NULL
)
-- +goose StatementEnd

INSERT INTO cartStatus (status) VALUES ('active'), ('abandoned'), ('completed');

-- +goose Down
-- +goose StatementBegin
DROP TABLE cartStatus;
-- +goose StatementEnd
