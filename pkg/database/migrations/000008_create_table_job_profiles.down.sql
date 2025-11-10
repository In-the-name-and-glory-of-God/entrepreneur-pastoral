-- Triggers must be dropped before the table.
DROP TRIGGER IF EXISTS set_timestamp_job_profiles ON job_profiles;
DROP TABLE IF EXISTS job_profile_fields_of_work;
DROP TABLE IF EXISTS job_profiles;
