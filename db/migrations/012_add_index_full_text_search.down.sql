DROP INDEX idx_pharmacy_products;

DROP INDEX product_search_vector_idx;

ALTER TABLE "products" DROP COLUMN "product_ts";

ALTER TABLE "products" DROP COLUMN "product_agg";