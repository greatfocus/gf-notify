CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_channel_staticName ON channel USING BTREE(staticName);