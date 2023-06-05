CREATE TABLE "toplists" (
  "id" serial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" varchar,
  "user_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "list_items" (
  "id" serial PRIMARY KEY,
  "toplist_id" int NOT NULL,
  "rank" int NOT NULL,
  "title" varchar NOT NULL,
  "description" varchar
);

CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "list_items" ("toplist_id");

ALTER TABLE "toplists" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "list_items" ADD FOREIGN KEY ("toplist_id") REFERENCES "toplists" ("id") ON DELETE CASCADE;

ALTER TABLE list_items ADD CONSTRAINT unique_rank_per_list UNIQUE (toplist_id, rank);

ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);