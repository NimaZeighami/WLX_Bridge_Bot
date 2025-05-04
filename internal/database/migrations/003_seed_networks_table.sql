-- +goose Up
INSERT INTO networks (id, name, symbol) VALUES 
    (1, 'Ethereum', 'ETH'),
    (56, 'Binance Smart Chain', 'BNB'),
    (728126428 , 'Tron', 'TRX'),
    (137, 'Polygon', 'POL');

-- +goose Down
DELETE FROM networks WHERE id IN (1, 56, 728126428, 137);