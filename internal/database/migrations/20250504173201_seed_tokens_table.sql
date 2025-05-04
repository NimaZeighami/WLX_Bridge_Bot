-- +goose Up

INSERT INTO tokens (symbol, name, network_id, contract_address, decimals, is_native) VALUES
    -- Ethereum
    ('ETH', 'Ethereum', 1, NULL , 18, true),
    ('USDT', 'Tether USD', 1, '0xdac17f958d2ee523a2206206994597c13d831ec7', 6, false),

    -- BSC 
    ('BNB', 'BNB', 56, NULL , 18, true),
    ('USDT', 'Tether USD', 56, '0x55d398326f99059ff775485246999027b3197955', 18, false),

    -- Polygon
    ('POL', 'Polygon', 137, NULL , 18, true),
    ('USDT', 'Tether USD', 137, '0xc2132d05d31c914a87c6611c10748aeb04b58e8f', 6, false),

    -- Tron
    ('TRX', 'Tron', 728126428, NULL , 6, true),
    ('USDT', 'Tether USD', 728126428, 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t', 6, false);

-- +goose Down
DELETE FROM tokens;
ALTER TABLE tokens AUTO_INCREMENT = 1;

