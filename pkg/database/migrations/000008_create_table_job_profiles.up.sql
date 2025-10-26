-- Table: job_profiles
-- Stores user-specific professional profile data (one-to-one with users).
CREATE TABLE IF NOT EXISTS job_profiles (
    user_id UUID PRIMARY KEY, -- This is both PK and FK
    open_to_work BOOLEAN NOT NULL DEFAULT FALSE,
    cv_path TEXT, -- Stores a URL or file system path to the CV

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE -- If the user is deleted, delete their profile too
        ON UPDATE CASCADE
);

-- Apply the trigger to 'updated_at' column
CREATE TRIGGER set_timestamp_job_profiles
BEFORE UPDATE ON job_profiles
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Table: job_profile_fields_of_work
-- Junction table for the many-to-many relationship between job_profiles and fields_of_work.
CREATE TABLE IF NOT EXISTS job_profile_fields_of_work (
    user_id UUID NOT NULL,
    field_of_work_id SMALLINT NOT NULL,

    -- Constraints
    PRIMARY KEY (user_id, field_of_work_id), -- Composite primary key
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES job_profiles(user_id) -- Links to job_profiles table
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_field_of_work
        FOREIGN KEY(field_of_work_id)
        REFERENCES fields_of_work(id)
        ON DELETE RESTRICT -- Prevents deleting a field if it's in use
        ON UPDATE CASCADE
);
