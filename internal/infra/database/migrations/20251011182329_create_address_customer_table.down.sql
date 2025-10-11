ALTER TABLE addresses ADD COLUMN customer_id uuid NOT NULL REFERENCES customers (id);
DROP TABLE address_customers;
