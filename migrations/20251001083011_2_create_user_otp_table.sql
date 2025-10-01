-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_otp_code()
    RETURNS TRIGGER AS $$
        DECLARE
            chars TEXT := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
            result TEXT := '';
            i INT := 0;
        BEGIN
            FOR i IN 1..6 LOOP
                result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
            END LOOP;

            NEW.otp_code := result;
            RETURN NEW;
        END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user_otp (
    user_id     INTEGER PRIMARY KEY,
    otp_code    CHAR(6) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at  TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '10 minutes'),
    is_used     BOOLEAN NOT NULL DEFAULT FALSE,
    attempts    SMALLINT NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_generate_otp_code
    BEFORE INSERT ON user_otp
        FOR EACH ROW
            EXECUTE FUNCTION generate_otp_code();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_generate_otp_code ON user_otp;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS user_otp CASCADE;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS generate_otp_code();
-- +goose StatementEnd
