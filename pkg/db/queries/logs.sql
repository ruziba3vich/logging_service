-- name: InsertLog :exec
INSERT INTO logs (message, event_time, level, service)
VALUES (?, ?, ?, ?);

-- name: InsertErrorLog :exec
INSERT INTO logs (message, event_time, level, service)
VALUES (?, ?, ?, ?);
