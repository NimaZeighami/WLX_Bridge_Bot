-- +goose Up
CREATE TABLE IF NOT EXISTS networks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    chain_id  INT UNIQUE NOT NULL              
);

-- +goose Down
DROP TABLE IF EXISTS networks;
