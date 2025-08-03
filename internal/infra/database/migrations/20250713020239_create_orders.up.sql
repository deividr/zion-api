CREATE TABLE orders (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_number int NOT NULL,
    pickup_date timestamp NOT NULL,
    customer_id uuid NOT NULL REFERENCES customers (id),
    employee_id text NOT NULL,
    order_local text,
    observations text,
    is_picked_up bool DEFAULT false,
    is_deleted boolean DEFAULT false NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    updated_at timestamp
);

CREATE TABLE order_products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id uuid NOT NULL REFERENCES orders (id),
    product_id uuid NOT NULL REFERENCES products (id),
    quantity integer NOT NULL,
    unity_type char(2) NOT NULL,
    price integer NOT NULL
);

CREATE TABLE order_sub_products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_product_id uuid NOT NULL REFERENCES order_products (id),
    product_id uuid NOT NULL REFERENCES products (id)
);
