BEGIN;

CREATE TABLE IF NOT EXISTS scheduled_event (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  start TIMESTAMP NOT NULL,
  end TIMESTAMP NOT NULL,
  title VARCHAR(40),
)

CREATE TABLE IF NOT EXISTS tutor_session (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  tutor_id UUID REFERENCES client (id),
  updated_by UUID REFERENCES client (id),
  room_id VARCHAR(40) NOT NULL,
  state VARCHAR(15) NOT NULL,
  event_id INT REFERENCES scheduled_event (id)
)

CREATE TABLE IF NOT EXISTS session_client (
  id SERIAL PRIMARY KEY,
  client_id UUID REFERENCES client (id),
  session_id INT REFERENCES tutor_session (id),
)

CREATE TABLE IF NOT EXISTS client_selected_session (
  id SERIAL PRIMARY KEY,
  can_attend BOOLEAN NOT NULL,
  client_id UUID REFERENCES client (id),
  session_id INT REFERENCES tutor_session (id),
  UNIQUE (session_id, client_id)
)

CREATE TABLE IF NOT EXISTS client_event (
  id SERIAL PRIMARY KEY,
  client_id UUID REFERENCES client (id),
  event_id INT REFERENCES scheduled_event (id)
  UNIQUE (ref_type, ref_id, page)
)

COMMIT;