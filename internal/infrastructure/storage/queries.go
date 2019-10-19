package storage

const (
	getUserQuery = `
SELECT id, username, role
FROM "user"
WHERE id = $1`
	getUserByCredentialsQuery = `
SELECT id, username, role
FROM "user"
WHERE username = $1
  AND password = $2`
)
