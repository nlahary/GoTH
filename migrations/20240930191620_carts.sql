-- +goose Up
-- +goose StatementBegin
CREATE TABLE carts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (status_id) REFERENCES cartStatus(status)
);
-- +goose StatementEnd

CREATE TRIGGER update_timestamp
AFTER UPDATE ON carts
FOR EACH ROW
BEGIN
    UPDATE carts 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE id = OLD.id;
END;

CREATE INDEX idx_carts_user_id ON carts (user_id);
CREATE INDEX idx_carts_status_id ON carts (status_id);

-- +goose Down
-- +goose StatementBegin
DROP TABLE carts;
-- +goose StatementEnd
