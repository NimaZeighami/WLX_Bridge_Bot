-- +goose Up
CREATE TABLE quotes (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    from_token_address VARCHAR(100) NOT NULL,
    to_token_address VARCHAR(100) NOT NULL,
    from_chain VARCHAR(32) NOT NULL,
    to_chain VARCHAR(32) NOT NULL,
    from_address VARCHAR(100) NOT NULL,
    to_address VARCHAR(100) NOT NULL,
    from_amount DECIMAL(65, 0) NOT NULL,
    to_amount_min DECIMAL(65, 0) NOT NULL,
    tx_hash VARCHAR(100),
    state ENUM('pending', 'submitted', 'confirmed', 'failed', 'expired', 'success') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- States for 

-- pending
-- → Quote has been created, but no transaction yet.

-- submitted (optional)
-- → The transaction is constructed and submitted to the blockchain.

-- confirmed
-- → The transaction is mined/confirmed on-chain.

-- failed
-- → Transaction failed (e.g., out of gas, user rejected, on-chain error).

-- expired
-- → Quote was not used within its valid window (often quotes are valid for a short time like 30s).

-- success
-- → Swap succeeded, funds bridged successfully.



-- +goose Down
DROP TABLE IF EXISTS quotes;
