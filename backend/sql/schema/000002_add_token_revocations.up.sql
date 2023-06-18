CREATE TABLE "token_revocations" (
  "token" varchar UNIQUE NOT NULL,
  "revoked_at" timestamptz NOT NULL DEFAULT (now())
);