
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    creation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users(username);


CREATE TABLE IF NOT EXISTS providers (
  id bigserial PRIMARY KEY,
  name varchar NOT NULL,
  slug varchar NOT NULL,
  status int NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

