ALTER TABLE addresses DROP CONSTRAINT fk_addresses_customer_id;
ALTER TABLE customers DROP CONSTRAINT unique_phone2;
ALTER TABLE customers DROP CONSTRAINT unique_phone;
ALTER TABLE products ALTER COLUMN name TYPE varchar(50);
