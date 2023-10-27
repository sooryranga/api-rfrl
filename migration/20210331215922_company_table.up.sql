BEGIN;

CREATE TABLE IF NOT EXISTS company (
  company_name VARCHAR(40) PRIMARY KEY,
  photo VARCHAR(120) NOT NULL,
  industry VARCHAR(40),
  about TEXT,
  active BOOLEAN NOT NULL DEFAULT FALSE,
);

CREATE TABLE IF NOT EXISTS company_email (
  email_domain VARCHAR(40) PRIMARY KEY,
  company_name VARCHAR(40) REFERENCES company (company_name),
  suggestions INT NOT NULL DEFAULT 0
  active BOOLEAN NOT NULL DEFAULT FALSE
)

ALTER TABLE IF EXISTS client
  ADD COLUMN company_name VARCHAR(40);

ALTER TABLE IF EXISTS client
  DROP CONSTRAINT IF EXISTS fk_company_name;

ALTER TABLE client
  ADD CONSTRAINT fk_company_name
  FORIEGN KEY (company_name)
  REFERENCES company(company_name)

COMMIT;