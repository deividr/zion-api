ALTER TABLE addresses DROP COLUMN is_default;
ALTER TABLE address_customers ADD COLUMN is_default boolean NOT NULL DEFAULT false;
