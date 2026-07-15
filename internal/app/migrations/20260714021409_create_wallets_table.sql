-- +goose Up
CREATE TABLE wallets (
    id         UUID PRIMARY KEY,
    balance    BIGINT NOT NULL,
    CONSTRAINT chk_balance_positive CHECK (balance >= 0)
);

-- +goose Down
DROP TABLE IF EXISTS wallets;
