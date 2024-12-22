CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
  id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name       varchar(50) NOT NULL,
  value      integer NOT NULL,
  unity_type  char(2) NOT NULL,
  is_deleted  boolean DEFAULT false NOT NULL,
  created_at  timestamp DEFAULT now() NOT NULL,
  updated_at  timestamp
);
