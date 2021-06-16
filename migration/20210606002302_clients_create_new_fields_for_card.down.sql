BEGIN;

ALTER TABLE client DROP COLUMN linkedin_profile;

ALTER TABLE client DROP COLUMN github_profile;

ALTER TABLE client DROP COLUMN years_of_experience;

ALTER TABLE client DROP COLUMN work_title;

COMMIT;