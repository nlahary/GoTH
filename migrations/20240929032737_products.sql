-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    desc TEXT NOT NULL,
    price DOUBLE NOT NULL
);
-- +goose StatementEnd

INSERT INTO products (name, desc, price) VALUES ('Product 1', 'Description 1', 100.00);
INSERT INTO products (name, desc, price) VALUES ('Product 2', 'Description 2', 200.00);
INSERT INTO products (name, desc, price) VALUES ('Product 3', 'Description 3', 300.00);
INSERT INTO products (name, desc, price) VALUES ('Product 4', 'Description 4', 400.00);
INSERT INTO products (name, desc, price) VALUES ('Product 5', 'Description 5', 500.00);
INSERT INTO products (name, desc, price) VALUES ('Product 6', 'Description 6', 600.00);


-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
