CREATE TABLE users
(
  id BIGINT NOT NULL,
  segment VARCHAR,
  create_at TIMESTEMP NOT NULL,
  delete_at TIMESTEMP
);

CREATE TABLE segments
(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL
);
