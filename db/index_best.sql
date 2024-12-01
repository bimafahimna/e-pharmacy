CREATE INDEX idx_pharmacy_products ON pharmacy_products(product_id, stock, price, sold_amount desc nulls last);

-- CREATE INDEX idx_products_id ON products(id);

CREATE INDEX idx_pharmacies ON pharmacies(is_active);

CREATE INDEX idx_partners ON partners(is_active);

-- DROP
DROP INDEX IF EXISTS idx_pharmacy_products;

DROP INDEX IF EXISTS idx_products_id;

DROP INDEX IF EXISTS idx_pharmacies;

DROP INDEX IF EXISTS idx_partners;