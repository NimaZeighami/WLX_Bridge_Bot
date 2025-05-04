-- +goose Up
CREATE TABLE IF NOT EXISTS tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    network_id INT NOT NULL,
    contract_address VARCHAR(300) DEFAULT '0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee',
    decimals INT NOT NULL,
    is_native BOOLEAN DEFAULT FALSE,


    INDEX idx_network_id (network_id),
    FOREIGN KEY (network_id) REFERENCES networks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS tokens;
