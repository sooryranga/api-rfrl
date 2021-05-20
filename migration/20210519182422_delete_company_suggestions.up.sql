BEGIN;

ALTER TABLE company_email 
DROP COLUMN suggestions;

COMMIT;