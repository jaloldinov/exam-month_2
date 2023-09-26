CREATE TABLE "branch" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "address" varchar,
  "phone_number" varchar,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

CREATE TABLE "category" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "parent_id" uuid,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

CREATE TABLE "product" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar UNIQUE NOT NULL,
  "category_id" uuid,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

CREATE TYPE coming_status AS ENUM ('in_process', 'finished');

CREATE TABLE "coming_table" (
  "id" uuid PRIMARY KEY,
  "coming_id" varchar NOT NULL,
  "branch_id" uuid,
  "date_time" timestamp,
  "status" coming_status DEFAULT 'in_process',
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

CREATE TABLE "coming_table_product" (
  "id" uuid PRIMARY KEY,
  "category_id" uuid,
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar  NOT NULL,
  "count" numeric NOT NULL DEFAULT 0,
  "total_price" numeric DEFAULT 0,
  "coming_table_id" uuid,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

CREATE TABLE "remaining" (
  "id" uuid PRIMARY KEY,
  "branch_id" uuid,
  "category_id" uuid,
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar  NOT NULL,
  "count" numeric NOT NULL DEFAULT 0,
  "total_price" numeric DEFAULT 0,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

ALTER TABLE "category" ADD FOREIGN KEY ("parent_id") REFERENCES "category" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "coming_table" ADD FOREIGN KEY ("branch_id") REFERENCES "branch" ("id");

ALTER TABLE "coming_table_product" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "coming_table_product" ADD FOREIGN KEY ("coming_table_id") REFERENCES "coming_table" ("id");

ALTER TABLE "remaining" ADD FOREIGN KEY ("branch_id") REFERENCES "branch" ("id");

ALTER TABLE "remaining" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");