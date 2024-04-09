CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL
);

CREATE TABLE global_chat
(
    id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users(id),
    message TEXT NOT NULL
);

CREATE TABLE private_chats
(
    id SERIAL PRIMARY KEY,
    sender_id INTEGER REFERENCES users(id),
    recipient_id INTEGER REFERENCES users(id),
    message TEXT NOT NULL
);
