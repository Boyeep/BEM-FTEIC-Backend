ALTER TABLE galeri ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'all';
UPDATE galeri SET category = CASE
    WHEN title ILIKE '%elektro%' THEN 'teknik_elektro'
    WHEN title ILIKE '%informatika%' THEN 'teknik_informatika'
    WHEN title ILIKE '%sistem informasi%' THEN 'sistem_informasi'
    WHEN title ILIKE '%biomedik%' THEN 'teknik_biomedik'
    WHEN title ILIKE '%teknologi informasi%' THEN 'teknologi_informasi'
    WHEN title ILIKE '%komputer%' THEN 'teknik_komputer'
    ELSE category
END WHERE category = 'all';

CREATE TABLE IF NOT EXISTS signup_whitelist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    created_by UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS site_visitors (
    id TEXT PRIMARY KEY,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_path TEXT,
    user_agent TEXT
);
ALTER TABLE site_visitors ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE site_visitors ADD COLUMN IF NOT EXISTS last_path TEXT;
ALTER TABLE site_visitors ADD COLUMN IF NOT EXISTS user_agent TEXT;

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
        COALESCE(NULLIF(NEW.raw_user_meta_data->>'username', ''), split_part(COALESCE(NEW.email, ''), '@', 1)),
        'member'
    )
    ON CONFLICT (id) DO NOTHING;
    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;
CREATE TRIGGER on_auth_user_created
AFTER INSERT ON auth.users
FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();

CREATE OR REPLACE FUNCTION public.is_admin()
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = public
AS $$
    SELECT EXISTS (
        SELECT 1 FROM profiles
        WHERE id = auth.uid() AND role = 'admin'
    );
$$;

CREATE OR REPLACE FUNCTION public.is_signup_email_whitelisted(candidate_email TEXT)
RETURNS BOOLEAN
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = public
AS $$
    SELECT EXISTS (
        SELECT 1 FROM signup_whitelist
        WHERE lower(email) = lower(trim(candidate_email))
    );
$$;

CREATE OR REPLACE FUNCTION public.get_visitor_count()
RETURNS BIGINT
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = public
AS $$ SELECT count(*) FROM site_visitors; $$;

-- Supabase Before User Created hook. Enable public.hook_validate_signup in
-- Authentication > Hooks after applying this migration.
CREATE OR REPLACE FUNCTION public.hook_validate_signup(event JSONB)
RETURNS JSONB
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
DECLARE
    candidate TEXT := lower(trim(event->'user'->>'email'));
BEGIN
    IF NOT public.is_signup_email_whitelisted(candidate) THEN
        RETURN jsonb_build_object(
            'error', jsonb_build_object(
                'http_code', 403,
                'message', 'Email is not authorized to sign up'
            )
        );
    END IF;
    RETURN '{}'::jsonb;
END;
$$;

ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE blogs ENABLE ROW LEVEL SECURITY;
ALTER TABLE events ENABLE ROW LEVEL SECURITY;
ALTER TABLE galeri ENABLE ROW LEVEL SECURITY;
ALTER TABLE signup_whitelist ENABLE ROW LEVEL SECURITY;
ALTER TABLE site_visitors ENABLE ROW LEVEL SECURITY;

DO $$
DECLARE policy_record RECORD;
BEGIN
    FOR policy_record IN
        SELECT schemaname, tablename, policyname
        FROM pg_policies
        WHERE schemaname = 'public'
          AND tablename IN ('profiles','blogs','events','galeri','signup_whitelist','site_visitors')
    LOOP
        EXECUTE format('DROP POLICY IF EXISTS %I ON %I.%I',
            policy_record.policyname, policy_record.schemaname, policy_record.tablename);
    END LOOP;
END $$;

DROP POLICY IF EXISTS profiles_public_read ON profiles;
CREATE POLICY profiles_public_read ON profiles FOR SELECT
USING (true);
DROP POLICY IF EXISTS profiles_self_update ON profiles;
CREATE POLICY profiles_self_update ON profiles FOR UPDATE TO authenticated
USING (id = auth.uid()) WITH CHECK (id = auth.uid());
DROP POLICY IF EXISTS profiles_self_insert ON profiles;
CREATE POLICY profiles_self_insert ON profiles FOR INSERT TO authenticated
WITH CHECK (id = auth.uid() AND role = 'member');

DROP POLICY IF EXISTS blogs_public_read ON blogs;
CREATE POLICY blogs_public_read ON blogs FOR SELECT
USING (status = 'PUBLISHED' OR public.is_admin());
DROP POLICY IF EXISTS blogs_admin_write ON blogs;
CREATE POLICY blogs_admin_write ON blogs FOR ALL TO authenticated
USING (public.is_admin()) WITH CHECK (public.is_admin());

DROP POLICY IF EXISTS events_public_read ON events;
CREATE POLICY events_public_read ON events FOR SELECT
USING (status = 'PUBLISHED' OR public.is_admin());
DROP POLICY IF EXISTS events_admin_write ON events;
CREATE POLICY events_admin_write ON events FOR ALL TO authenticated
USING (public.is_admin()) WITH CHECK (public.is_admin());

DROP POLICY IF EXISTS galeri_public_read ON galeri;
CREATE POLICY galeri_public_read ON galeri FOR SELECT USING (true);
DROP POLICY IF EXISTS galeri_admin_write ON galeri;
CREATE POLICY galeri_admin_write ON galeri FOR ALL TO authenticated
USING (public.is_admin()) WITH CHECK (public.is_admin());

DROP POLICY IF EXISTS whitelist_admin_all ON signup_whitelist;
CREATE POLICY whitelist_admin_all ON signup_whitelist FOR ALL TO authenticated
USING (public.is_admin()) WITH CHECK (public.is_admin());

DROP POLICY IF EXISTS visitors_insert ON site_visitors;
CREATE POLICY visitors_insert ON site_visitors FOR INSERT TO anon, authenticated
WITH CHECK (true);
DROP POLICY IF EXISTS visitors_update ON site_visitors;
CREATE POLICY visitors_update ON site_visitors FOR UPDATE TO anon, authenticated
USING (true) WITH CHECK (true);

REVOKE ALL ON FUNCTION public.is_admin() FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.is_admin() TO authenticated;
GRANT EXECUTE ON FUNCTION public.is_signup_email_whitelisted(TEXT) TO anon, authenticated;
GRANT EXECUTE ON FUNCTION public.get_visitor_count() TO anon, authenticated;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'supabase_auth_admin') THEN
        GRANT EXECUTE ON FUNCTION public.hook_validate_signup(JSONB) TO supabase_auth_admin;
    END IF;
END $$;

REVOKE ALL ON profiles, blogs, events, galeri, signup_whitelist, site_visitors FROM anon, authenticated;
GRANT SELECT (id, username, avatar_url, updated_at) ON profiles TO anon, authenticated;
GRANT SELECT ON blogs, events, galeri TO anon, authenticated;
GRANT INSERT (id, last_seen_at, last_path, user_agent) ON site_visitors TO anon, authenticated;
GRANT UPDATE (last_seen_at, last_path, user_agent) ON site_visitors TO anon, authenticated;
