-- +goose Up
CREATE TABLE IF NOT EXISTS networks (
    id INT PRIMARY KEY ,              
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(50) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS networks;
