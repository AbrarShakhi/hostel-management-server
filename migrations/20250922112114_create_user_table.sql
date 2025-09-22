-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
    uuid            SERIAL          PRIMARY KEY,
    email           VARCHAR(128)    UNIQUE,
    phone           VARCHAR(16)     UNIQUE NOT NULL,
    password        VARCHAR(256),
    first_name      VARCHAR(128),
    last_name       VARCHAR(64),
    gardian_name    VARCHAR(128),   NOT NULL,
    gardian_phone   VARCHAR(16),    NOT NULL,
    created_on      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
-- +goose StatementEnd
