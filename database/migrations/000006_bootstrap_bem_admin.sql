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
