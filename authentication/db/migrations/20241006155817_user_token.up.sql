CREATE TABLE user_token
(
    id UUID NOT NULL PRIMARY KEY,
    userId UUID NOT NULL,
    token VARCHAR(255),
    CONSTRAINT user_id_fk FOREIGN KEY (userId) REFERENCES auth_user (id) ON DELETE CASCADE
);