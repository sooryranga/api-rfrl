BEGIN;

ALTER TABLE IF EXISTS client
  ADD COLUMN institution VARCHAR(40);

ALTER TABLE IF EXISTS client
  ADD COLUMN degree VARCHAR(40);

ALTER TABLE IF EXISTS client
  ADD COLUMN field_of_study VARCHAR(40);

ALTER TABLE IF EXISTS client
  ADD COLUMN start_year SMALLINT;

ALTER TABLE IF EXISTS client
  ADD COLUMN end_year SMALLINT;

END;