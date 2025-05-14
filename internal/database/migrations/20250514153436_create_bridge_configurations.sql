-- +goose Up
CREATE TABLE IF NOT EXISTS bridge_configurations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    network VARCHAR(255) NOT NULL,
    chain_id INT,
    token VARCHAR(255) NOT NULL,
    token_contract_address VARCHAR(255) NOT NULL,
    token_decimals INT NOT NULL,
    bridgers_smart_contract_address VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at VARCHAR(255),
    updated_at VARCHAR(255)
);

-- +goose Down
DROP TABLE IF EXISTS bridge_configurations;


-- TODO: This table is temporary and will be removed in the future 