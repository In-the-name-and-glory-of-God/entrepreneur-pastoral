-- Table: industries
-- Stores the different business industries. The 'key' field contains translation keys.
CREATE TABLE IF NOT EXISTS industries (
    id SMALLSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE
);

-- Add default industries with translation keys
INSERT INTO industries (key)
VALUES
    ('industry.technology'),
    ('industry.healthcare'),
    ('industry.finance'),
    ('industry.education'),
    ('industry.retail'),
    ('industry.manufacturing'),
    ('industry.construction'),
    ('industry.hospitality'),
    ('industry.transportation'),
    ('industry.real_estate'),
    ('industry.agriculture'),
    ('industry.media_entertainment'),
    ('industry.telecommunic ations'),
    ('industry.energy'),
    ('industry.food_beverage'),
    ('industry.professional_services'),
    ('industry.nonprofit'),
    ('industry.government'),
    ('industry.other')
ON CONFLICT (key) DO NOTHING;