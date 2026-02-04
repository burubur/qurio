ALTER TABLE sources
ADD COLUMN sync_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE sources
ADD COLUMN sync_schedule TEXT DEFAULT 'daily';
-- minute, hourly, daily
ALTER TABLE sources
ADD COLUMN last_synced_at TIMESTAMP WITH TIME ZONE;