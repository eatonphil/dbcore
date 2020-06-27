CREATE TABLE users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME,
);

CREATE TABLE notes (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  note TEXT NOT NULL,
  created_by BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME,
  FOREIGN KEY (created_by) REFERENCES users (id)
);
