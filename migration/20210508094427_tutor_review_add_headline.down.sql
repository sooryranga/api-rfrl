BEGIN;

ALTER TABLE IF EXISTS tutor_review
  DROP COLUMN IF EXISTS headline;

DROP TABLE IF EXISTS pending_tutor_review;

COMMIT;