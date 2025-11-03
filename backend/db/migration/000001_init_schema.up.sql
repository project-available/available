CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "role" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "student_id" varchar UNIQUE NOT NULL,
  "is_delete" bool DEFAULT false NOT NULL
);

CREATE TABLE "rooms" (
  "id" bigserial PRIMARY KEY,
  "location" varchar NOT NULL,
  "name" varchar NOT NULL,
  "image" varchar NOT NULL,
  "is_delete" bool DEFAULT false NOT NULL
);

CREATE TABLE "bookings" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "room_id" bigint NOT NULL,
  "start" timestamptz NOT NULL,
  "end" timestamptz NOT NULL,
  "status" varchar DEFAULT 'pending' NOT NULL,
  "phone_booking" varchar NOT NULL
);

CREATE TABLE "custom_fields" (
  "id" bigserial PRIMARY KEY,
  "key" varchar UNIQUE NOT NULL,
  "shown" bool NOT NULL DEFAULT true
);

CREATE TABLE "custom_fields_value" (
  "id" bigserial PRIMARY KEY,
  "room_id" bigint NOT NULL,
  "customfield_id" bigint NOT NULL,
  "value" varchar NOT NULL
);

CREATE INDEX ON "accounts" ("name");

CREATE INDEX ON "rooms" ("name");

CREATE UNIQUE INDEX ON "custom_fields_value" ("room_id", "customfield_id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "bookings" ADD CONSTRAINT bookings_start_before_end CHECK ("start" < "end");

ALTER TABLE "custom_fields_value" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE;

ALTER TABLE "custom_fields_value" ADD FOREIGN KEY ("customfield_id") REFERENCES "custom_fields" ("id") ON DELETE CASCADE;