-- +goose Up
CREATE TABLE quotes (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    from_chain VARCHAR(32) NOT NULL,
    from_token_address VARCHAR(100) NOT NULL,
    to_chain VARCHAR(32) NOT NULL,
    to_token_address VARCHAR(100) NOT NULL,
    from_address VARCHAR(100) NOT NULL,
    to_address VARCHAR(100) NOT NULL,
    from_amount DECIMAL(65, 0) NOT NULL,
    to_amount_min DECIMAL(65, 0) NOT NULL,
    tx_hash VARCHAR(100) ,
    state ENUM('started', 'approved', 'broadcast', 'verified', 'approval_failed', 'broadcast_failed', 'failed') NOT NULL DEFAULT 'started',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- +goose Down
DROP TABLE IF EXISTS quotes;
