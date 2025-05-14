-- +goose Up
CREATE TABLE IF NOT EXISTS network_token_pairs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    source_network_symbol VARCHAR(255) NOT NULL,
    source_token_symbol VARCHAR(255) NOT NULL,
    target_network_symbol VARCHAR(255) NOT NULL,
    target_token_symbol VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE KEY unique_pair (source_network, source_token, target_network, target_token)
);

--  Chain sypmbol is based on The Bridgers cross-chain bridge

-- +goose Down
DROP TABLE IF EXISTS network_token_pairs;