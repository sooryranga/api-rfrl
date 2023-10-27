BEGIN;

CREATE TABLE IF NOT EXISTS tags (
  id SERIAL PRIMARY KEY,
  tag_name VARCHAR(40) UNIQUE,
  about TEXT
);

CREATE TABLE IF NOT EXISTS question (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  title VARCHAR(120) NOT NULL,
  body TEXT NOT NULL,
  applicants INT DEFAULT 0,
  images text[],
  from_id UUID REFERENCES client (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question_applicants (
  question_id SERIAL REFERENCES question (id) ON DELETE CASCADE,
  applicant_id UUID REFERENCES client (id) ON DELETE CASCADE,
  UNIQUE (question_id, applicant_id)
);

CREATE TABLE IF NOT EXISTS question_tags (
  id SERIAL PRIMARY KEY,
  question_id SERIAL REFERENCES question (id) ON DELETE CASCADE,
  tag_id SERIAL REFERENCES tags (id) ON DELETE CASCADE,
  UNIQUE (tag_id, question_id)
);

COMMIT;
