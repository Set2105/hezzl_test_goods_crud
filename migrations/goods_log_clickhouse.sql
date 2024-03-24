CREATE TABLE IF NOT EXISTS nats_goods_log_queue (
    id Int8 NOT NULL,
    project_id Int8 NOT NULL,
    "name" String NOT NULL,
    description String NOT NULL DEFAULT '',
    priority Int8,
    removed Bool NOT NULL DEFAULT false,
    event_time DateTime NOT NULL,
  ) ENGINE = NATS
    SETTINGS nats_url = 'nats:4222',
             nats_subjects = 'goods_log',
             nats_format = 'JSONEachRow';

CREATE TABLE IF NOT EXISTS  goods_log (
	id Int8 NOT NULL,
	project_id Int8 NOT NULL,
	"name" String NOT NULL,
	description String NOT NULL DEFAULT '',
	priority Int8,
	removed Bool NOT NULL DEFAULT false,
	event_time DateTime NOT NULL,
) 
ENGINE = MergeTree
PRIMARY KEY (id, project_id,"name");

CREATE MATERIALIZED VIEW IF NOT EXISTS  nats_goods_log_consumer TO goods_log
AS SELECT id, project_id, "name", description, priority, removed, event_time FROM nats_goods_log_queue;