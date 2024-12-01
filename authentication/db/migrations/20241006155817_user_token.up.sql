CREATE TABLE user_token
(
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    user_ip VARCHAR (255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_in BIGINT NOT NULL,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES auth_user (id) ON DELETE CASCADE
);