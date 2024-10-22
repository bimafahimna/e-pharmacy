ALTER TABLE "products" ADD COLUMN product_agg VARCHAR;

ALTER TABLE "products" ADD COLUMN product_ts tsvector;

CREATE INDEX product_search_vector_idx ON products USING gin(product_ts);

CREATE INDEX idx_pharmacy_products ON pharmacy_products(product_id, stock, price, sold_amount desc nulls last);
