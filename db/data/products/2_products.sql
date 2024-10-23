UPDATE "products" SET "product_agg" = concat("name",' ',"generic_name",' ',"categories",' ', "description", ' ', "selling_unit");

UPDATE	"products" SET "product_ts" = to_tsvector("product_agg");