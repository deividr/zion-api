CREATE TABLE customers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    phone text NOT NULL UNIQUE,
    phone2 text UNIQUE,
    email text,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp,
    old_id int
);
