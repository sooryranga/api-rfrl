BEGIN;

ALTER TABLE IF EXISTS auth 
  DROP COLUMN IF EXISTS sign_up_flow;

COMMIT;