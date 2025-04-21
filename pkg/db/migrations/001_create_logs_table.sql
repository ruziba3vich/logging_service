CREATE TABLE IF NOT EXISTS logs (
	id UUID DEFAULT generateUUIDv4(),
	message String,
	event_time DateTime,
	level String,
	service String,
	received_at DateTime DEFAULT now()
) ENGINE = MergeTree
PARTITION BY toYYYYMM(event_time)
ORDER BY (event_time)
TTL event_time + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;
