-- Existing admins become part of the whitelist authorization source.
INSERT INTO signup_whitelist (email)
SELECT LOWER(BTRIM(email))
FROM profiles
WHERE role = 'admin' AND BTRIM(email) <> ''
ON CONFLICT (email) DO NOTHING;

-- Existing accounts that were already whitelisted receive dashboard access.
UPDATE profiles AS profile
SET role = 'admin', updated_at = NOW()
FROM signup_whitelist AS whitelist
WHERE LOWER(BTRIM(profile.email)) = LOWER(BTRIM(whitelist.email))
  AND profile.role <> 'admin';

-- New accounts inherit admin access from the whitelist used by the signup hook.
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
    INSERT INTO profiles (id, email, username, role)
    VALUES (
        NEW.id,
        COALESCE(NEW.email, ''),
        COALESCE(
            NULLIF(NEW.raw_user_meta_data->>'username', ''),
            split_part(COALESCE(NEW.email, ''), '@', 1)
        ),
        CASE
            WHEN public.is_signup_email_whitelisted(COALESCE(NEW.email, '')) THEN 'admin'
            ELSE 'member'
        END
    )
    ON CONFLICT (id) DO UPDATE
    SET
        email = EXCLUDED.email,
        username = EXCLUDED.username,
        role = EXCLUDED.role,
        updated_at = NOW();
    RETURN NEW;
END;
$$;
