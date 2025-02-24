USE connectfour;

CREATE TABLE user
(
    id    BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
    email VARCHAR(255)                   NOT NULL,
    name  VARCHAR(255)                   NOT NULL,
    token VARCHAR(255)                   NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_user_email ON user (email);


CREATE TABLE game
(
    game_key       VARCHAR(20) NOT NULL,
    player1_id     BIGINT      NOT NULL,
    player2_id     BIGINT      NULL,
    created_at     DATETIME    NOT NULL,
    started_at     DATETIME    NULL,
    finished_at    DATETIME    NOT NULL,
    player_turn_id BIGINT      NULL,
    public         bool        NOT NULL,
    status         VARCHAR(20) NOT NULL,
    board_json     TEXT        NULL,
    PRIMARY KEY (game_key)
)

