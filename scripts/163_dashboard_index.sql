CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_dashboard_report ON dashboard USING BTREE(year, month);