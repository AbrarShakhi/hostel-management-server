-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_ (
    id            SERIAL       PRIMARY KEY,
    email         VARCHAR(128) UNIQUE NOT NULL,
    phone         VARCHAR(20)  UNIQUE NOT NULL,
    password_     VARCHAR(256) NOT NULL,
    first_name    VARCHAR(128) NOT NULL,
    last_name     VARCHAR(64)  NOT NULL,
    date_of_birth DATE         NOT NULL,
    gender        CHAR(1)      NOT NULL,
    nationality   VARCHAR(16)  NOT NULL,
    created_on    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login    TIMESTAMP
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE user_address (
    id             SERIAL       PRIMARY KEY,
    street_address VARCHAR(256) NOT NULL,
    city           VARCHAR(128) NOT NULL,
    state_         VARCHAR(128),
    postal_code    VARCHAR(16),
    country        VARCHAR(128) NOT NULL,
    user_id        INTEGER      NOT NULL, 
    FOREIGN KEY (user_id) REFERENCES user_(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE emergency_contact (
    id           SERIAL       PRIMARY KEY,
    name_        VARCHAR(128) NOT NULL,
    phone        VARCHAR(20)  NOT NULL,
    relationship VARCHAR(64),
    user_id      INTEGER      NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user_(id) ON DELETE CASCADE
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS emergency_contact CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS user_address CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS user_ CASCADE;
-- +goose StatementEnd
