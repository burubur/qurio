-- Drop the old constraint
ALTER TABLE sources DROP CONSTRAINT IF EXISTS sources_content_hash_key;

-- Create a partial unique index that excludes deleted records
CREATE UNIQUE INDEX sources_content_hash_active_idx ON sources (content_hash) WHERE deleted_at IS NULL;
