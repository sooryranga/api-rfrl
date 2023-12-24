BEGIN;

CREATE TABLE client_wanting_company_referral (
  company_id INT REFERENCES company (id),
  client_id UUID REFERENCES client (id),
  UNIQUE (company_id, client_id)
);

ALTER TABLE client ADD COLUMN is_looking_for_referral BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;