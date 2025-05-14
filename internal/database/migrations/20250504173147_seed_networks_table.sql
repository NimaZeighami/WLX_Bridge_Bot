-- +goose Up
INSERT INTO networks (chain_id, name, symbol) VALUES 
    (1, 'Ethereum', 'ETH'),
    (56, 'Binance Smart Chain', 'BSC'),
    (728126428 , 'Tron', 'TRX'),
    (137, 'Polygon', 'POLYGON'),
 ;
 
--  Chain sypmbol is based on The Bridgers cross-chain bridge

-- +goose Down
DELETE FROM networks;
ALTER TABLE networks AUTO_INCREMENT = 1;
