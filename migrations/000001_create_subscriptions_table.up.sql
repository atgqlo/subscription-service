CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK(price >= 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL CHECK (start_date ~ '^[0-1][0-9]-[0-9]{4}$'),
    end_date VARCHAR(7),
    created_at TIMESTAMPTZ DEFAULT NOW()
);