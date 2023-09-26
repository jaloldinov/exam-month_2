
-- Drop foreign key constraints
ALTER TABLE "remaining" DROP CONSTRAINT IF EXISTS "remaining_category_id_fkey";
ALTER TABLE "remaining" DROP CONSTRAINT IF EXISTS "remaining_branch_id_fkey";
ALTER TABLE "coming_table_product" DROP CONSTRAINT IF EXISTS "coming_table_product_coming_table_id_fkey";
ALTER TABLE "coming_table_product" DROP CONSTRAINT IF EXISTS "coming_table_product_category_id_fkey";
ALTER TABLE "coming_table" DROP CONSTRAINT IF EXISTS "coming_table_branch_id_fkey";
ALTER TABLE "product" DROP CONSTRAINT IF EXISTS "product_category_id_fkey";
ALTER TABLE "category" DROP CONSTRAINT IF EXISTS "category_parent_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "remaining";
DROP TABLE IF EXISTS "coming_table_product";
DROP TABLE IF EXISTS "coming_table";
DROP TABLE IF EXISTS "product";
DROP TABLE IF EXISTS "category";
DROP TABLE IF EXISTS "branch";

DROP TYPE IF EXISTS coming_status;
