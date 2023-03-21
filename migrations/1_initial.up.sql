BEGIN;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  telegram_user_nickname text,
  telegram_user_name text,
  telegram_user_id BIGINT NOT NULL UNIQUE,
  telegram_chat_id BIGINT NOT NULL UNIQUE,
  user_status INT DEFAULT 0
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  telegram_user_id INT NOT NULL,
  telegram_chat_id INT NOT NULL,
  name text,
  weight REAL,
  alreadyused BOOLEAN DEFAULT false,
  inlist BOOLEAN DEFAULT false,
  infridge BOOLEAN DEFAULT false,
  intrash BOOLEAN DEFAULT false,
  created_at TIMESTAMP WITH TIME ZONE,
  finished_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  last_update TIMESTAMP WITH TIME ZONE,
  rest_time BIGINT DEFAULT 0,
  timerenable BOOLEAN DEFAULT false
);


COMMIT;