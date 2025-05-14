-- +goose Up
INSERT INTO networks (chain_id, name, symbol) VALUES 
    (1, 'Ethereum', 'ETH'),
    (56, 'Binance Smart Chain', 'BNB'),
    (728126428 , 'Tron', 'TRX'),
    (137, 'Polygon', 'POL');

-- +goose Down
DELETE FROM networks;
ALTER TABLE networks AUTO_INCREMENT = 1;
