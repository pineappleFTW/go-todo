--  psql -f ./pkg/models/postgres/teardown.sql -d todo

DROP TABLE todos;
DROP TABLE refresh_tokens;
DROP TABLE users cascade;