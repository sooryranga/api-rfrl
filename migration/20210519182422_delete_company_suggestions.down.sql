BEGIN;

ALTER TABLE company_email 
ADD COLUMN suggestions INT DEFAULT 0;

COMMIT;