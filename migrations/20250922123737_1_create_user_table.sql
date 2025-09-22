-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_address (
    id             SERIAL       PRIMARY KEY,
    street_address VARCHAR(256) NOT NULL,
    city           VARCHAR(128) NOT NULL,
    state_         VARCHAR(128),
    postal_code    VARCHAR(16),
    country        VARCHAR(128) NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE emergency_contact (
    id           SERIAL       PRIMARY KEY,
    name_        VARCHAR(128) NOT NULL,
    phone        VARCHAR(20)  NOT NULL,
    relationship VARCHAR(64)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE user_ (
    id            SERIAL       PRIMARY KEY,
    email         VARCHAR(128) UNIQUE NOT NULL,
    phone         VARCHAR(20)  UNIQUE NOT NULL,
    password_     VARCHAR(256) NOT NULL,
    first_name    VARCHAR(128) NOT NULL,
    last_name     VARCHAR(64)  NOT NULL,
    date_of_birth DATE         NOT NULL,
    gender        CHAR(1)      NOT NULL CHECK (gender IN ('M', 'F', 'O')),
    nationality   VARCHAR(16)  NOT NULL,
    created_on    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login    TIMESTAMP,
    address_      INTEGER      NOT NULL,
    emergency_    INTEGER      NOT NULL,
    FOREIGN KEY (address_) REFERENCES user_address(id) ON DELETE SET NULL,
    FOREIGN KEY (emergency_) REFERENCES emergency_contact(id) ON DELETE SET NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_ CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS user_address CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS emergency_contact CASCADE;
-- +goose StatementEnd
