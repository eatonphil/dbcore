CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL
);

CREATE TABLE notes (
  id BIGSERIAL PRIMARY KEY,
  note TEXT NOT NULL,
  created_by BIGINT NOT NULL,
  FOREIGN KEY (created_by) REFERENCES users (id)
);
