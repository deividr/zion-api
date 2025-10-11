CREATE TABLE address_customers (
    address_id uuid NOT NULL REFERENCES addresses (id),
    customer_id uuid NOT NULL REFERENCES customers (id),
    PRIMARY KEY (address_id, customer_id)
);

ALTER TABLE addresses DROP COLUMN customer_id;
