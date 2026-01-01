CREATE TABLE products (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  old_id     integer,
  name       text NOT NULL,
  value      integer NOT NULL,
  unity_type  char(2) NOT NULL,
  is_deleted  boolean DEFAULT false NOT NULL,
  created_at  timestamp DEFAULT now() NOT NULL,
  updated_at  timestamp
);
