CREATE TABLE addresses (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    old_id int,
    customer_id uuid NOT NULL,
    cep text NOT NULL,
    street text NOT NULL,
    adr_number text NOT NULL,
    neighborhood text NOT NULL,
    city text NOT NULL,
    adr_state text NOT NULL,
    aditional_details text NOT NULL,
    distance integer NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp
);
