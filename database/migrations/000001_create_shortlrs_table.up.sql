BEGIN;
CREATE TABLE IF NOT EXISTS "shortlrs"(
    "id" UUID PRIMARY KEY,
    "long_url" VARCHAR(255) NOT NULL,
    "short_url" VARCHAR(255) NOT NULL,
    "access_count" BIGINT DEFAULT 0,
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP
);
COMMIT;