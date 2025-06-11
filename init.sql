CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    transaction_id VARCHAR(255) UNIQUE NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    state VARCHAR(10) NOT NULL CHECK (state IN ('win', 'lose')),
    source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('game', 'server', 'payment')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_transaction_id ON transactions(transaction_id);

INSERT INTO users (id, balance) VALUES 
(1, 0.00),
(2, 0.00),
(3, 0.00);