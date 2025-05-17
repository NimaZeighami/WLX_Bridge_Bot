-- +goose Up
CREATE TABLE IF NOT EXISTS network_token_pairs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_network_symbol VARCHAR(255) NOT NULL,
    from_token_symbol VARCHAR(255) NOT NULL,
    to_network_symbol VARCHAR(255) NOT NULL,
    to_token_symbol VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS network_token_pairs;