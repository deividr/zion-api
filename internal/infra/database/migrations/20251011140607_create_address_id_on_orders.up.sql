ALTER TABLE orders ADD COLUMN address_id uuid REFERENCES addresses (id);
