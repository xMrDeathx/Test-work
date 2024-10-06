CREATE TABLE auth_user
(
    id         UUID         NOT NULL PRIMARY KEY,
    email       VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);