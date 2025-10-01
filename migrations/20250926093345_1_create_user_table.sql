-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_ (
    id            SERIAL       PRIMARY KEY,
    email         VARCHAR(128) UNIQUE NOT NULL  CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}(?:\.[A-Za-z]{2,})*$'),
    phone         VARCHAR(30)  UNIQUE NOT NULL,
    password_     VARCHAR(256) ,
    first_name    VARCHAR(128) NOT NULL         CHECK (first_name ~* '^[A-Za-z .-]+$'),
    last_name     VARCHAR(64)                   CHECK (last_name IS NULL OR last_name ~* '^[A-Za-z .-]+$'),
    date_of_birth DATE         NOT NULL         CHECK (date_of_birth <= CURRENT_DATE),
    gender        CHAR(1)      NOT NULL         CHECK (gender IN ('M', 'F', 'O')),
    nationality   VARCHAR(16)  NOT NULL,
    created_on    TIMESTAMPTZ  NOT NULL         DEFAULT CURRENT_TIMESTAMP,
    last_login    TIMESTAMPTZ  ,
    has_left      BOOLEAN      NOT NULL         DEFAULT FALSE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE user_address (
    id             SERIAL       PRIMARY KEY,
    street_address VARCHAR(256) NOT NULL,
    city           VARCHAR(128) NOT NULL         CHECK (city ~* '^[A-Za-z\s\/\.\-]+$'),
    state_         VARCHAR(128)                  CHECK (state_ IS NULL OR state_ ~* '^[A-Za-z\s\/\.\-]+$'),
    postal_code    VARCHAR(16)                   CHECK (postal_code IS NULL OR postal_code ~* '^[A-Za-z0-9\s-]+$'),
    country        VARCHAR(128) NOT NULL         CHECK (country ~* '^[A-Za-z\s\/\.\-]+$'),
    FOREIGN KEY (id) REFERENCES user_(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE emergency_contact (
    id           SERIAL       PRIMARY KEY,
    name_        VARCHAR(128) NOT NULL,
    phone        VARCHAR(30)  NOT NULL,
    relationship VARCHAR(64)  NOT NULL      CHECK (relationship ~* '^[A-Za-z .-]+$'),
    FOREIGN KEY (id) REFERENCES user_(id) ON DELETE CASCADE
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_address CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS emergency_contact CASCADE;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS user_ CASCADE;
-- +goose StatementEnd
