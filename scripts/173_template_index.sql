CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_template_active ON template USING BTREE(deleted);