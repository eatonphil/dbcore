INSERT INTO organizations
VALUES (1, 'Admin', DATETIME('now'), DATETIME('now'), null);

-- password: admin
INSERT INTO users (is_admin, organization, username, password, name, created_at, updated_at)
VALUES (TRUE, 1, 'admin', '$2y$12$rH9QTzZJmPGIPIfofMfRsOh8vD5u612VYOduOvq951vIp7ddQn4Ai', 'Admin', DATETIME('now'), DATETIME('now'));

INSERT INTO organizations
VALUES (2, 'Notes Today', DATETIME('now'), DATETIME('now'), null);

-- password: admin
INSERT INTO users (is_admin, organization, username, password, name, created_at, updated_at)
VALUES (TRUE, 2, 'notes-admin', '$2y$12$rH9QTzZJmPGIPIfofMfRsOh8vD5u612VYOduOvq951vIp7ddQn4Ai', 'Notes Admin', DATETIME('now'), DATETIME('now'));

-- password: editor
INSERT INTO users (is_admin, organization, username, password, name, created_at, updated_at)
VALUES (FALSE, 3, 'editor', '$2y$12$F5RlC/VhW1LNcCZi5ZZ1P.JnBvBShzSiMt3rKqZjhvJvylT6bJu2i', 'Editor', DATETIME('now'), DATETIME('now'));
