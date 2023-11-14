BEGIN;

ALTER TABLE IF EXISTS tutor_session
  DROP COLUMN IF EXISTS conference_id;

END;