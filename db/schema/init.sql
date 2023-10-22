CREATE EXTENSION pg_trgm;

CREATE TABLE categories(
  id serial primary key,
  reg text unique,
  group_id int
);

CREATE TABLE category_groups(
  id serial primary key,
  name varchar(128),
  severity int
);

CREATE TABLE logs_data(
  id text,
  ts timestamptz,
  data text,
  category_id int,
  resolved_at timestamptz
);

CREATE INDEX categories_reg_idx ON categories USING gin(reg gin_trgm_ops);
