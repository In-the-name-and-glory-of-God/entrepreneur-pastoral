CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign Key to the Business table
    business_id UUID NOT NULL,

    -- Product Details
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    -- Price (Decimal type for accurate money handling)
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0), 

    -- Images
    image_url VARCHAR(255),

    -- Status and Availability
    is_available BOOLEAN NOT NULL DEFAULT TRUE,

    -- Timestamps (Good practice, often useful for sorting and caching)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_business
        FOREIGN KEY(business_id)
        REFERENCES business(id)
        ON DELETE RESTRICT -- Prevents deleting a business if products are still assigned to it
        ON UPDATE CASCADE
);