CREATE TABLE users
(
  id BIGINT NOT NULL,
  segment VARCHAR,
  create_at TIMESTAMP NOT NULL,
  delete_at TIMESTAMP
);

CREATE TABLE segments
(
  id SERIAL PRIMARY KEY,
  name VARCHAR NOT NULL
);
