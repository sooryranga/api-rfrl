BEGIN;

CREATE TABLE IF NOT EXISTS session_event (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  start TIMESTAMP NOT NULL,
  end TIMESTAMP NOT NULL,
  title VARCHAR(40)
)

CREATE TABLE IF NOT EXISTS tutor_session (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  tutor_id UUID REFERENCES client (id),
  updated_by UUID REFERENCES client (id),
  room_id VARCHAR(40) NOT NULL,
  state VARCHAR(15) NOT NULL,
  target_event_id INT REFERENCES event (id)
)

ALTER TABLE IF EXISTS session_event 
  ADD COLUMN session_id INT;

ALTER TABLE session_event 
   ADD CONSTRAINT fk_session_id
   FOREIGN KEY (session_id) 
   REFERENCES tutor_session ( id );


CREATE TABLE IF NOT EXISTS session_client (
  id SERIAL PRIMARY KEY,
  client_id UUID REFERENCES client (id),
  session_id INT REFERENCES tutor_session (id),
)

CREATE TABLE IF NOT EXISTS client_selected_event (
  id SERIAL PRIMARY KEY,
  client_id UUID REFERENCES client (id),
  session_id INT REFERENCES tutor_session (id),
  event_id INT REFERENCES session_event (id)
)

COMMIT;