BEGIN;

ALTER TABLE IF EXISTS tutor_review
  ADD COLUMN headline VARCHAR(120);

CREATE TABLE IF NOT EXISTS pending_tutor_review (
  mentee_id UUID REFERENCES client (id) ON DELETE CASCADE,
  tutor_id UUID REFERENCES client (id) ON DELETE CASCADE,
  PRIMARY KEY (mentee_id, tutor_id)
);

COMMIT;
