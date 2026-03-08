ALTER TABLE users
    ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN last_login_at TIMESTAMP WITH TIME ZONE NULL; -- Can switch to NOT NULL in future migration if needed
