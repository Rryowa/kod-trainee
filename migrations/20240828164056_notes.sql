-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     username TEXT UNIQUE NOT NULL,
     password TEXT NOT NULL
);
CREATE INDEX users_name ON users USING hash(username);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL
);
CREATE INDEX notes_user_id ON notes USING hash(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_name;
DROP INDEX IF EXISTS notes_user_id;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd