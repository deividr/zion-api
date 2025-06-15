CREATE TABLE IF NOT EXISTS category_products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		description TEXT
);

INSERT INTO category_products (name, description) VALUES ('Diversos', 'Diversos produtos');
INSERT INTO category_products (name, description) VALUES ('Massas', 'Massas e diversos derivados da farinha de trigo');
INSERT INTO category_products (name, description) VALUES ('Bebidas', 'Bebidas e liquidos');
INSERT INTO category_products (name, description) VALUES ('Carnes', 'Carnes, frangos, peixes, etc.');
INSERT INTO category_products (name, description) VALUES ('Saladas', 'Maionese, salpicão, antepastos, etc.');
INSERT INTO category_products (name, description) VALUES ('Sobremesas', 'Sobremesas, doces, etc.');