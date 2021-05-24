BEGIN;

ALTER TABLE client
  DROP CONSTRAINT IF EXISTS fk_company_name;

ALTER TABLE IF EXISTS client
  DROP COLUMN company_name;

ALTER TABLE client
  ADD COLUMN company_id INT;

ALTER TABLE IF EXISTS company_email
  DROP CONSTRAINT IF EXISTS company_email_company_name_fkey;

ALTER TABLE IF EXISTS company_email
  DROP COLUMN company_name;

ALTER TABLE IF EXISTS company_email
  ADD COLUMN company_id INT;

ALTER TABLE company DROP CONSTRAINT company_pkey;

ALTER TABLE company ADD UNIQUE (company_name);

ALTER TABLE company ADD COLUMN id SERIAL;

ALTER TABLE company ADD PRIMARY KEY (id);

ALTER TABLE company_email
  ADD CONSTRAINT company_email_company_id_fkey
  FOREIGN KEY (company_id)
  REFERENCES company(id) 
  ON DELETE CASCADE;

ALTER TABLE client
  ADD CONSTRAINT fk_company_id
  FOREIGN KEY (company_id)
  REFERENCES company(id) 
  ON DELETE NO ACTION;

COMMIT;
