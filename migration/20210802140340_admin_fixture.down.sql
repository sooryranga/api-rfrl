BEGIN;

DELETE FROM auth
	WHERE client_id = 'af496484-7c7c-45a8-a409-96d61351f43a';

DELETE FROM client 
  WHERE id ='af496484-7c7c-45a8-a409-96d61351f43a';

END;