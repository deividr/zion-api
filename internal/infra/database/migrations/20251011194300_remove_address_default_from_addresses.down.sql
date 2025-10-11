ALTER TABLE addresses ADD COLUMN is_default bool DEFAULT false NOT NULL;
ALTER TABLE address_customers DROP COLUMN is_default;
