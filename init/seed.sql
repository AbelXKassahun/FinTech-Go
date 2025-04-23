CREATE TABLE fee_configs (
    id SERIAL PRIMARY KEY,
    transaction_type VARCHAR(50) NOT NULL,
    tier VARCHAR(50), -- NULLable for base rules applying to all tiers
    base_percentage NUMERIC(10, 5) NOT NULL DEFAULT 0, -- e.g., 0.015 for 1.5%
    min_fee NUMERIC(12, 2) NOT NULL DEFAULT 0,
    max_fee NUMERIC(12, 2), -- NULLable if no max fee
    peak_start_time TIME, -- Format HH:MM, e.g., '18:00'
    peak_end_time TIME,
    peak_surcharge NUMERIC(10, 5) NOT NULL DEFAULT 0, -- e.g., 0.005 for 0.5% surcharge
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (transaction_type, tier) -- Ensures only one rule per type/tier combo (handle NULL tier correctly)
);

CREATE UNIQUE INDEX unique_transaction_type_tier
ON fee_configs (transaction_type)
WHERE tier IS NULL;

CREATE INDEX idx_fee_configs_type_tier ON fee_configs (transaction_type, tier);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_fee_configs_updated_at
BEFORE UPDATE ON fee_configs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();




CREATE TABLE users (
    user_id VARCHAR(100) PRIMARY KEY,
    full_name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20),
    password_hash TEXT NOT NULL,-- use bcrypt hashes
    tier VARCHAR(50) NOT NULL DEFAULT 'Basic',
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER trigger_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Sample Data Insertion:
INSERT INTO fee_configs (transaction_type, tier, base_percentage, min_fee, max_fee, peak_start_time, peak_end_time, peak_surcharge) VALUES
('BILL_PAYMENT', 'Basic', 0.03, 1.00, 100.00, '18:00', '22:00', 0.005), -- Basic: 3% + 0.5% peak surcharge, max 100
('BILL_PAYMENT', 'Premium', 0.015, 0.50, 100.00, NULL, NULL, 0),   -- Premium: 1.5%, max 100, no peak surcharge [cite: 16]
('TRANSFER', 'Enterprise', 0.005, 0.10, NULL, NULL, NULL, 0),     -- Enterprise Transfer: 0.5%, no max
('WALLET_TOPUP', NULL, 0.01, 0.20, 50.00, NULL, NULL, 0);         -- Base rule for TopUp: 1%, max 50 (applies to all tiers unless overridden)

INSERT INTO users (user_id, full_name, email, phone_number, password_hash, tier, last_login) 
VALUES
('user123', 'Alice Johnson', 'alice.johnson@example.com', '+1234567890', 
    '$2a$10$0hBN26eojcn0euNg7NBBKuv.GUK0MlBqlxzAdLrJDbHFIV8x5piaK', 'Basic', '2025-04-22 10:00:00+00'),
('user456', 'Bob Smith', 'bob.smith@example.com', '+0987654321', 
    '$2a$10$bt.sKBhlphHhTTChBrdi7uZWeX2ibf41ZE857l95mD/MxaQ03CwW6', 'Premium', '2025-04-21 15:30:00+00'),
('user789', 'Charlie Davis', 'charlie.davis@example.com', '+1122334455', 
    '$2a$10$wsE02Soi6Z40h6xFTlDY9O2ZDF8tmPP5ij4TjS0dRYKvS7avu17di', 'Enterprise', NULL);
