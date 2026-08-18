-- Older production databases may already contain these tables without the
-- constraints declared by CREATE TABLE IF NOT EXISTS in earlier migrations.
-- Repair duplicate whitelist rows before adding the unique index required by
-- the ON CONFLICT clauses below and in subsequent migrations.
DELETE FROM signup_whitelist
WHERE ctid IN (
    SELECT ctid
    FROM (
        SELECT
            ctid,
            ROW_NUMBER() OVER (
                PARTITION BY LOWER(BTRIM(email))
                ORDER BY created_at ASC NULLS LAST, ctid ASC
            ) AS duplicate_number
        FROM signup_whitelist
    ) AS ranked_whitelist
    WHERE duplicate_number > 1
);

UPDATE signup_whitelist
SET email = LOWER(BTRIM(email))
WHERE email <> LOWER(BTRIM(email));

CREATE UNIQUE INDEX IF NOT EXISTS signup_whitelist_email_key
    ON signup_whitelist (email);

INSERT INTO profiles (id, email, username, role)
SELECT
    id,
    email,
    COALESCE(
        NULLIF(BTRIM(raw_user_meta_data ->> 'username'), ''),
        'Admin BEM FTEIC'
    ),
    'admin'
FROM auth.users
WHERE LOWER(BTRIM(email)) = 'website.fteic@gmail.com'
ON CONFLICT (id) DO UPDATE
SET
    email = EXCLUDED.email,
    role = 'admin',
    updated_at = NOW();

INSERT INTO signup_whitelist (email)
VALUES ('website.fteic@gmail.com')
ON CONFLICT (email) DO NOTHING;
