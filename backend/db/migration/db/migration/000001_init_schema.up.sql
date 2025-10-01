CREATE TYPE "booking_status" AS ENUM (
  'pending',
  'accepted',
  'canceled'
);

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "role" varchar,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar,
  "phone" varchar UNIQUE,
  "student_id" varchar UNIQUE NOT NULL,
  "is_delete" bool DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "update_at" timestamptz DEFAULT (now()),
  "update_by" varchar
);

CREATE TABLE "rooms" (
  "id" bigserial PRIMARY KEY,
  "capacity" int,
  "location" varchar,
  "name" varchar,
  "is_delete" bool DEFAULT false,
  "created_at" timestamptz DEFAULT (now()),
  "update_at" timestamptz DEFAULT (now()),
  "update_by" varchar
);

CREATE TABLE "bookings" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "room_id" bigint NOT NULL,
  "start" timestamptz NOT NULL,
  "end" timestamptz NOT NULL,
  "status" booking_status DEFAULT 'pending',
  "phone_booking" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "update_at" timestamptz DEFAULT (now()),
  "update_by" varchar
);

CREATE INDEX ON "accounts" ("name");

CREATE INDEX ON "rooms" ("name");

COMMENT ON COLUMN "rooms"."capacity" IS 'can only be positive';

COMMENT ON COLUMN "rooms"."location" IS 'H1-100';

ALTER TABLE "bookings" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");
