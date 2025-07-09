ALTER TABLE products ALTER COLUMN name TYPE text;
ALTER TABLE customers ADD CONSTRAINT unique_phone UNIQUE (phone);
ALTER TABLE customers ADD CONSTRAINT unique_phone2 UNIQUE (phone2);
ALTER TABLE addresses ADD CONSTRAINT fk_addresses_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id);
