-- Define the possible values for the job type
CREATE TYPE job_type_enum AS ENUM (
    'Full-Time',
    'Part-Time',
    'Contract'
);

-- Define the possible values for the job location (Remote or On Site)
CREATE TYPE job_location_enum AS ENUM (
    'Remote',
    'On Site',
    'Hybrid' -- Added 'Hybrid' as a common third option
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign Key to the Business table
    business_id UUID NOT NULL,

    -- Job Details
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    
    -- Controlled Values using ENUM types
    type job_type_enum NOT NULL,
    location job_location_enum NOT NULL,
    
    -- Application link
    application_link VARCHAR(255),
    
    -- Status
    is_open BOOLEAN NOT NULL DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_business
        FOREIGN KEY(business_id)
        REFERENCES business(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);