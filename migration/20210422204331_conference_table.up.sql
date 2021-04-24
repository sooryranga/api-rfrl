BEGIN;

CREATE TABLE IF NOT EXISTS conference_code (
  id SERIAL PRIMARY KEY,
  code TEXT NOT NULL,
  result TEXT,
);

CREATE TABLE IF NOT EXISTS session_conference (
  session_id SERIAL REFERENCES tutor_session ( id ) ON DELETE CASCADE,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  code_state VARCHAR(15) NOT NULL DEFAULT 'not_running',
  latest_code SERIAL REFERENCES conference_code ( id ) ON DELETE RESTRICT
);

COMMIT;