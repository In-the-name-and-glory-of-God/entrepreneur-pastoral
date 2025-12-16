CREATE TABLE church (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    diocese VARCHAR(255) NOT NULL,
    parish_number VARCHAR(50),
    website_url VARCHAR(255),
    phone_number VARCHAR(20),
    address_id UUID NOT NULL,
    is_archdiocese BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- Foreign Key Constraint linking to the Address table
    CONSTRAINT fk_address
        FOREIGN KEY (address_id)
        REFERENCES address(id)
        ON DELETE RESTRICT -- Prevents deleting an address if a church is still linked to it
);