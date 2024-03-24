CREATE SCHEMA IF NOT EXISTS test;
CREATE TABLE IF NOT EXISTS  test.projects(
	id serial8 NOT NULL,
	"name" varchar(256) NOT NULL,
	created_at timestamp DEFAULT current_timestamp NOT NULL,
	CONSTRAINT projects_id_pk PRIMARY KEY (id)
);
COMMENT ON COLUMN test.projects.id IS 'id записи';
COMMENT ON COLUMN test.projects.name IS 'название';
COMMENT ON COLUMN test.projects.created_at IS 'дата и время';
INSERT INTO test.projects (id, "name") VALUES (1, 'первая запись') ON CONFLICT (id) DO NOTHING;
CREATE TABLE IF NOT EXISTS test.goods (
	id serial8,
	project_id int8 NOT NULL,
	"name" varchar(256) NOT NULL,
	description text NOT NULL DEFAULT '',
	priority serial8,
	removed bool NOT NULL DEFAULT false,
	created_at timestamp NOT NULL default current_timestamp,
	CONSTRAINT goods_pk PRIMARY KEY (id, project_id)
);
COMMENT ON COLUMN test.goods.id IS 'id записи';
COMMENT ON COLUMN test.goods.project_id IS 'id кампании';
COMMENT ON COLUMN test.goods.name IS 'название';
COMMENT ON COLUMN test.goods.description IS 'описание';
COMMENT ON COLUMN test.goods.priority IS 'приоритет';
COMMENT ON COLUMN test.goods.removed IS 'статус удаления';
COMMENT ON COLUMN test.goods.created_at IS 'дата и время';
CREATE INDEX IF NOT EXISTS goods_name_indx ON test.goods ("name");
