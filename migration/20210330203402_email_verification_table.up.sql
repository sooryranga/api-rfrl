BEGIN;

CREATE TABLE IF NOT EXISTS email_verification (
  id SERIAL PRIMARY KEY,
  client_id UUID REFERENCES client (id) ON DELETE CASCADE,
  email VARCHAR(40),
  email_type VARCHAR(20),
  pass_code VARCHAR(6),
  UNIQUE (client_id, email_type)
);

COMMIT;