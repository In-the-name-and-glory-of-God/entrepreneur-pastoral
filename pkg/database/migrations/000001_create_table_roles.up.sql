-- Table: roles
-- Stores the user roles available in the system.
CREATE TABLE IF NOT EXISTS roles (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

-- Add default roles
INSERT INTO roles (name, description)
VALUES
    ('Admin', 'Administrator with full system access'),
    ('Manager', 'User with full Admin module access'),
    ('Assistent', 'User with basic Admin module access'),
    ('Entrepreneur', 'Standard entrepreneur with basic permissions'),
    ('User', 'Standard user with basic permissions')
ON CONFLICT (name) DO NOTHING;