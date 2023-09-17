CREATE TABLE pets (
  id bigserial not null primary key,
  name varchar,
  created_at timestamp,
  updated_at timestamp
);