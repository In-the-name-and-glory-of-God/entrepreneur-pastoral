CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign Key to the Business table
    business_id UUID NOT NULL,

    -- Service Details
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    
    -- Price (Decimal type for accurate money handling)
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_business
        FOREIGN KEY(business_id)
        REFERENCES business(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);