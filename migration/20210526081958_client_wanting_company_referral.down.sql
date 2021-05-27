BEGIN;

ALTER TABLE client DROP COLUMN is_looking_for_referral;

DROP TABLE client_wanting_company_referral;

COMMIT;