-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    network_id INT NOT NULL,
    contract_address VARCHAR(300),
    decimals INT NOT NULL,
    is_native BOOLEAN DEFAULT FALSE,


    UNIQUE (name, symbol, network_id),
    INDEX idx_network_id (network_id),
    FOREIGN KEY (network_id) REFERENCES networks(chain_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
