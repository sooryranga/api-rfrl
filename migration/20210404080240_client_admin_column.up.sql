BEGIN;

ALTER TABLE IF EXISTS client
  ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;