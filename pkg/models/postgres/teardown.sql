--  psql -f ./pkg/models/postgres/teardown.sql -d todo

DROP TABLE todos;
DROP TABLE users;
DROP TABLE refresh_tokens;