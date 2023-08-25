CREATE TABLE
    "segment" (
        "name" VARCHAR(255) NOT NULL UNIQUE,
        "amount" FLOAT,
        CONSTRAINT "segments_pk" PRIMARY KEY ("name")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "people" (
        "user_id" integer NOT NULL,
        CONSTRAINT "user_pk" PRIMARY KEY ("user_id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "user_segment" (
        "user_id" integer NOT NULL,
        "seg_name" VARCHAR(255) NOT NULL,
        "delete_time" TIMESTAMPTZ
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "log" (
        "user_id" integer NOT NULL,
        "seg_name" VARCHAR(255) NOT NULL,
        "operation" VARCHAR(255) NOT NULL,
        "op_time" TIMESTAMPTZ NOT NULL
    )
WITH (OIDS = FALSE);

ALTER TABLE "user_segment"
ADD
    CONSTRAINT "user_segment_fk0" FOREIGN KEY ("user_id") REFERENCES "people"("user_id") ON DELETE CASCADE;

ALTER TABLE "user_segment"
ADD
    CONSTRAINT "user_segment_fk1" FOREIGN KEY ("seg_name") REFERENCES "segment"("seg_name") ON DELETE CASCADE;