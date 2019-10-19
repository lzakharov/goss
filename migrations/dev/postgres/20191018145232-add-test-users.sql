-- +migrate Up

INSERT INTO "user" (id, username, password, role)
VALUES (1, 'admin', 'password', 'admin')
ON CONFLICT DO NOTHING;

INSERT INTO "user" (id, username, password, role)
VALUES (2, 'user1', 'password', 'user'),
       (3, 'user2', 'password', 'user'),
       (4, 'user3', 'password', 'user')
ON CONFLICT DO NOTHING;

-- +migrate Down

DELETE
FROM "user"
WHERE id BETWEEN 1 AND 4;
