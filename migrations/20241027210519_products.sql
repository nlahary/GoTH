-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    desc TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    price DOUBLE NOT NULL
);
-- +goose StatementEnd

INSERT INTO products (name, desc, quantity, price) VALUES ('Product 1', 'Description 1', 19, 100);
INSERT INTO products (name, desc, quantity, price) VALUES ('Product 2', 'Description 2', 20, 100);
INSERT INTO products (name, desc, quantity, price) VALUES ('Product 3', 'Description 3', 5, 100);
INSERT INTO products (name, desc, quantity, price) VALUES ('Product 4', 'Description 4', 4, 100);
INSERT INTO products (name, desc, quantity, price) VALUES ('Product 5', 'Description 5', 3, 100);
INSERT INTO products (name, desc, quantity, price) VALUES ('Product 6', 'Description 6', 1, 100);


-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
