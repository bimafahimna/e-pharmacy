CREATE INDEX idx_products_name ON products USING GIN (to_tsvector('english', name));
CREATE INDEX idx_products_generic_name ON products USING GIN (to_tsvector('english', generic_name));
CREATE INDEX idx_products_description ON products USING GIN (to_tsvector('english', description));
CREATE INDEX idx_products_selling_unit ON products USING GIN (to_tsvector('english', selling_unit));
CREATE INDEX idx_products_categories ON products USING GIN (categories);

CREATE INDEX idx_manufacturers_name ON manufacturers USING GIN (to_tsvector('english', name));

CREATE INDEX idx_product_classifications_name ON product_classifications USING GIN (to_tsvector('english', name));

CREATE INDEX idx_product_forms_name ON product_forms USING GIN (to_tsvector('english', name));

DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_generic_name;
DROP INDEX IF EXISTS idx_products_description;
DROP INDEX IF EXISTS idx_products_selling_unit;
DROP INDEX IF EXISTS idx_products_categories;

DROP INDEX IF EXISTS idx_manufacturers_name;

DROP INDEX IF EXISTS idx_product_classifications_name;

DROP INDEX IF EXISTS idx_product_forms_name;
