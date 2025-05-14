-- +goose Up
INSERT INTO bridge_configurations (network, token, chain_id, token_contract_address, token_decimals, bridgers_smart_contract_address, is_enabled) VALUES

-- BSC Configuration
("BSC", "USDT", 56, "0x55d398326f99059ff775485246999027b3197955", 18, "0xb685760ebd368a891f27ae547391f4e2a289895b", true),

-- TRON Configuration
("TRON", "USDT", 728126428, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", 6, "0xb685760ebd368a891f27ae547391f4e2a289895b", true),

-- POL Configuration
("POL", "USDT", 137, "0xc2132d05d31c914a87c6611c10748aeb04b58e8f", 6, "0xb685760ebd368a891f27ae547391f4e2a289895b", true);

-- +goose Down
DELETE FROM bridge_configurations;
ALTER TABLE bridge_configurations AUTO_INCREMENT = 1;


-- TODO: This table is temporary and will be removed in the future 