CREATE TABLE "users" (
  "chat_id" varchar PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "active" bool NOT NULL
);

CREATE TABLE "grades" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "name" varchar NOT NULL,
  "grade" varchar NOT NULL,
  "range" varchar NOT NULL,
  "percentage" varchar NOT NULL,
  "course_name" varchar NOT NULL
);

ALTER TABLE "grades" ADD FOREIGN KEY ("owner") REFERENCES "users" ("chat_id");

CREATE INDEX ON "grades" ("owner");

CREATE INDEX ON "grades" ("course_name");
