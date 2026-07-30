ALTER TABLE events
    ADD COLUMN IF NOT EXISTS publication_status TEXT NOT NULL DEFAULT 'PUBLISHED';

UPDATE events
SET status = CASE
    WHEN status IN ('UPCOMING', 'ONGOING', 'ENDED') THEN status
    WHEN event_date > CURRENT_DATE THEN 'UPCOMING'
    WHEN event_date = CURRENT_DATE THEN 'ONGOING'
    ELSE 'ENDED'
END;

UPDATE events
SET publication_status = CASE
    WHEN publication_status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED') THEN publication_status
    ELSE 'PUBLISHED'
END;

UPDATE blogs
SET status = 'DRAFT'
WHERE status NOT IN ('DRAFT', 'PUBLISHED', 'ARCHIVED');

ALTER TABLE blogs DROP CONSTRAINT IF EXISTS blogs_status_check;
ALTER TABLE blogs
    ADD CONSTRAINT blogs_status_check
    CHECK (status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED'));

ALTER TABLE events DROP CONSTRAINT IF EXISTS events_status_check;
ALTER TABLE events
    ADD CONSTRAINT events_status_check
    CHECK (status IN ('UPCOMING', 'ONGOING', 'ENDED'));

ALTER TABLE events DROP CONSTRAINT IF EXISTS events_publication_status_check;
ALTER TABLE events
    ADD CONSTRAINT events_publication_status_check
    CHECK (publication_status IN ('DRAFT', 'PUBLISHED', 'ARCHIVED'));

CREATE INDEX IF NOT EXISTS idx_events_publication_date
    ON events (publication_status, event_date DESC);

DROP POLICY IF EXISTS events_public_read ON public.events;
CREATE POLICY events_public_read ON public.events
FOR SELECT TO anon, authenticated
USING (publication_status = 'PUBLISHED' OR public.is_admin());
