BEGIN;

INSERT INTO client (id, first_name, last_name, about, email, is_tutor, is_admin)
	VALUES ('af496484-7c7c-45a8-a409-96d61351f43a', 'Support', 'rfrl', 'rfrl Support', 'admin@rfrl.ca', TRUE, TRUE);

INSERT INTO auth (email, password_hash, auth_type, client_id, sign_up_flow)
VALUES ('admin@rfrl.ca', '\x243261243130247447645276514b424d374c577a6e463647765655324f304c4f6676387075416a6d4336416f416c4a30733044534d647357746f4957', 'email', 'af496484-7c7c-45a8-a409-96d61351f43a', 100);

COMMIT;