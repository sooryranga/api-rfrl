BEGIN;

ALTER TABLE IF EXISTS auth
  ADD COLUMN sign_up_flow SMALLINT NOT NULL DEFAULT 0;

COMMIT;