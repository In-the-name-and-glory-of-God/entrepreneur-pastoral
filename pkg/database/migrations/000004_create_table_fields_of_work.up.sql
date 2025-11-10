-- Table: fields_of_work
-- Stores the different professional fields of work.
CREATE TABLE IF NOT EXISTS fields_of_work (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Add some default fields of work
INSERT INTO fields_of_work (name)
VALUES
    ('Architecture and Engineering Occupations'),
    ('Arts, Design, Entertainment, Sports, and Media Occupations'),
    ('Building and Grounds Cleaning and Maintenance Occupations'),
    ('Business and Financial Operations Occupations'),
    ('Community and Social Services Occupations'),
    ('Computer and Mathematical Occupations'),
    ('Construction and Extraction Occupations'),
    ('Education, Training, and Library Occupations'),
    ('Farming, Fishing, and Forestry Occupations'),
    ('Food Preparation and Serving Related Occupations'),
    ('Healthcare Practitioners and Technical Occupations'),
    ('Healthcare Support Occupations'),
    ('Installation, Maintenance, and Repair Occupations'),
    ('Legal Occupations'),
    ('Life, Physical, and Social Science Occupations'),
    ('Management Occupations'),
    ('Military Specific Occupations'),
    ('Office and Administrative Support Occupations'),
    ('Personal Care and Service Occupations'),
    ('Production Occupations'),
    ('Protective Service Occupations'),
    ('Sales and Related Occupations'),
    ('Other')
ON CONFLICT (name) DO NOTHING;