-- +goose Up
-- +goose StatementBegin
CREATE TABLE cartItems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cart_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
)
-- +goose StatementEnd

CREATE INDEX idx_cartItems_cart_id ON cartItems (cart_id);
CREATE INDEX idx_cartItems_product_id ON cartItems (product_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
