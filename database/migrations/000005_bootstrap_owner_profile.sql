INSERT INTO profiles (id, email, username, role)
SELECT
    id,
    email,
    COALESCE(
        NULLIF(BTRIM(raw_user_meta_data ->> 'username'), ''),
        SPLIT_PART(email, '@', 1)
    ),
    'admin'
FROM auth.users
WHERE LOWER(email) = 'stevenprobot@gmail.com'
ON CONFLICT (id) DO UPDATE
SET
    email = EXCLUDED.email,
    role = 'admin',
    updated_at = NOW();
