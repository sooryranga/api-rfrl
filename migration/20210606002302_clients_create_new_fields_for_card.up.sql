BEGIN;

ALTER TABLE client ADD COLUMN linkedin_profile VARCHAR(100);

ALTER TABLE client ADD COLUMN github_profile VARCHAR(100);

ALTER TABLE client ADD COLUMN years_of_experience SMALLINT;

ALTER TABLE client ADD COLUMN work_title VARCHAR(50);

COMMIT;