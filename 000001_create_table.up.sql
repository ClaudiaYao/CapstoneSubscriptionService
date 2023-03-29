
CREATE TABLE "subscription" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "playlist_id" varchar NOT NULL,
  "customized" bool NOT NULL,
  "status" varchar(20) NOT NULL,
  "frequency" varchar NOT NULL,
  "start_date" datetime NOT NULL,
  "end_date" datetime,
  "receiver_name" varchar NOT NULL,
  "receiver_contact" varchar NOT NULL
);

CREATE TABLE "dish_delivery" (
  "id" varchar PRIMARY KEY,
  "subscription_dish_id" varchar NOT NULL,
  "status" varchar(30) NOT NULL,
  "expected_time" datetime NOT NULL,
  "delivery_time" datetime,
  "note" varchar(100)
);

CREATE TABLE "subscription_dish" (
  "id" varchar PRIMARY KEY,
  "dish_id" varchar NOT NULL,
  "subscription_id" varchar NOT NULL,
  "schedule_time" datetime NOT NULL,
  "frequency" varchar NOT NULL,
  "dish_options" varchar NOT NULL,
  "note" varchar
);

CREATE INDEX ON "subscription" ("id");

CREATE INDEX ON "dish_delivery" ("id");

CREATE INDEX ON "subscription_dish" ("id");

ALTER TABLE "dish_delivery" ADD FOREIGN KEY ("subscription_dish_id") REFERENCES "subscription_dish" ("id");

ALTER TABLE "subscription_dish" ADD FOREIGN KEY ("subscription_id") REFERENCES "subscription" ("id");





