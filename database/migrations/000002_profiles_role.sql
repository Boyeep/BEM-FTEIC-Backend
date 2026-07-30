ALTER TABLE profiles ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'member';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'profiles_role_check') THEN
        ALTER TABLE profiles ADD CONSTRAINT profiles_role_check CHECK (role IN ('member', 'admin'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_profiles_role ON profiles (role);

-- Bootstrap the current project owner; subsequent role changes belong in Supabase.
UPDATE profiles SET role = 'admin' WHERE email = 'stevenprobot@gmail.com';
