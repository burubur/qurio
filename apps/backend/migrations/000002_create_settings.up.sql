CREATE TABLE IF NOT EXISTS settings (
    id INT PRIMARY KEY DEFAULT 1,
    rerank_provider TEXT NOT NULL DEFAULT 'none',
    rerank_api_key TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT settings_singleton CHECK (id = 1)
);

-- Insert default row if not exists
INSERT INTO settings (id, rerank_provider, rerank_api_key)
VALUES (1, 'none', '')
ON CONFLICT (id) DO NOTHING;
