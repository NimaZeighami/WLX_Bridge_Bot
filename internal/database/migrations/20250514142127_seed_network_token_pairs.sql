-- +goose Up

INSERT INTO network_token_pairs (from_network_symbol,from_token_symbol,to_network_symbol,to_token_symbol) VALUES

-- pol <--> bsc
("POLYGON", "USDT", "BSC", "USDT"),
("BSC", "USDT", "POLYGON", "USDT"),

-- trx <--> bsc
("TRX", "USDT", "BSC", "USDT"),
("BSC", "USDT", "TRX", "USDT"),

-- bsc <--> eth
("BSC", "USDT", "ETH", "USDT"),
("ETH", "USDT", "BSC", "USDT"),

-- eth <--> trx
("TRX", "USDT", "ETH", "USDT"),
("ETH", "USDT", "TRX", "USDT");


-- +goose Down
DELETE FROM network_token_pairs;
ALTER TABLE network_token_pairs AUTO_INCREMENT = 1;
