USE connectfour;

CREATE TABLE user
(
    id    INTEGER UNSIGNED AUTO_INCREMENT NOT NULL,
    email VARCHAR(255)                    NOT NULL,
    name  VARCHAR(255)                    NOT NULL,
    token VARCHAR(255)                    NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_user_email ON user (email);
