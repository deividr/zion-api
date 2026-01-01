CREATE TABLE addresses (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    old_id int,
    customer_id uuid NOT NULL REFERENCES customers (id),
    cep text NOT NULL,
    street text NOT NULL,
    number text NOT NULL,
    neighborhood text NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    aditional_details text NOT NULL,
    distance integer NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp,
    is_deleted bool DEFAULT false NOT NULL,
    is_default bool DEFAULT false NOT NULL
);
