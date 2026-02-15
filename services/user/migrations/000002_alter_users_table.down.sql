ALTER TABLE IF EXISTS users
    DROP COLUMN IF EXISTS email_verified,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS last_login_at;
