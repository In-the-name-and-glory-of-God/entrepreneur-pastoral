-- Table: fields_of_work
-- Stores the different professional fields of work. The 'key' field contains translation keys.
CREATE TABLE IF NOT EXISTS fields_of_work (
    id SMALLSERIAL PRIMARY KEY,
    key VARCHAR(255) NOT NULL UNIQUE
);

-- Add default fields of work with translation keys
INSERT INTO fields_of_work (key)
VALUES
    ('field_of_work.architecture_engineering'),
    ('field_of_work.arts_design_entertainment'),
    ('field_of_work.building_grounds'),
    ('field_of_work.business_financial'),
    ('field_of_work.community_social'),
    ('field_of_work.computer_mathematical'),
    ('field_of_work.construction_extraction'),
    ('field_of_work.education_training'),
    ('field_of_work.farming_fishing_forestry'),
    ('field_of_work.food_preparation'),
    ('field_of_work.healthcare_practitioners'),
    ('field_of_work.healthcare_support'),
    ('field_of_work.installation_maintenance'),
    ('field_of_work.legal'),
    ('field_of_work.life_physical_social_science'),
    ('field_of_work.management'),
    ('field_of_work.military'),
    ('field_of_work.office_administrative'),
    ('field_of_work.personal_care'),
    ('field_of_work.production'),
    ('field_of_work.protective_service'),
    ('field_of_work.sales'),
    ('field_of_work.other')
ON CONFLICT (key) DO NOTHING;