CREATE TABLE IF NOT EXISTS business (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign Keys
    user_id UUID NOT NULL REFERENCES users(id),
    industry_id SMALLINT NOT NULL REFERENCES industries(id),

    -- Core Data
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,

    -- Contact Information
    email CITEXT UNIQUE NOT NULL,
    phone_country_code VARCHAR(10),
    phone_number VARCHAR(20),
    website_url VARCHAR(255),
    logo_url VARCHAR(255),

    -- Status and Timestamps
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT -- Prevents deleting a user if businesses are still assigned to it
        ON UPDATE CASCADE,
    CONSTRAINT fk_industry
        FOREIGN KEY(industry_id)
        REFERENCES industries(id)
        ON DELETE RESTRICT -- Prevents deleting an industry if businesses are still assigned to it
        ON UPDATE CASCADE
);