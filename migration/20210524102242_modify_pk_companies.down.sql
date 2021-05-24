BEGIN;

ALTER TABLE company_email
  DROP CONSTRAINT company_email_company_id_fkey;

ALTER TABLE company_email
  DROP COLUMN company_id;

ALTER TABLE company_email
  ADD COLUMN company_name VARCHAR(40);

ALTER TABLE client
  DROP CONSTRAINT fk_company_id;

ALTER TABLE client
  DROP COLUMN company_id;

ALTER TABLE client
  ADD COLUMN company_name VARCHAR(40);

ALTER TABLE company DROP CONSTRAINT company_pkey;

ALTER TABLE company DROP COLUMN id;

ALTER TABLE company DROP CONSTRAINT company_company_name_key;

ALTER TABLE company ADD PRIMARY KEY (company_name);

ALTER TABLE company_email
  ADD CONSTRAINT company_email_company_name_fkey
  FOREIGN KEY (company_name)
  REFERENCES company(company_name) 
  ON DELETE CASCADE;

ALTER TABLE client
  ADD CONSTRAINT fk_company_name
  FOREIGN KEY (company_name)
  REFERENCES company(company_name) 
  ON DELETE NO ACTION;

COMMIT;
