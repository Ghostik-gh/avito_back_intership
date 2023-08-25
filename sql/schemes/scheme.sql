CREATE TABLE
    "segment" (
        "seg_id" serial NOT NULL,
        "name" VARCHAR(255) NOT NULL UNIQUE,
        "amount" FLOAT,
        CONSTRAINT "segments_pk" PRIMARY KEY ("seg_id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "user" (
        "user_id" serial NOT NULL,
        CONSTRAINT "user_pk" PRIMARY KEY ("user_id")
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "user_segment" (
        "user_id" integer NOT NULL,
        "seg_id" integer NOT NULL,
        "duration" TIMESTAMP
    )
WITH (OIDS = FALSE);

CREATE TABLE
    "log" (
        "user_id" integer NOT NULL,
        "seg_id" integer NOT NULL,
        "operation" VARCHAR(255) NOT NULL,
        "op_time" TIMESTAMP NOT NULL
    )
WITH (OIDS = FALSE);

ALTER TABLE "user_segment"
ADD
    CONSTRAINT "user_segment_fk0" FOREIGN KEY ("user_id") REFERENCES "user"("user_id");

ALTER TABLE "user_segment"
ADD
    CONSTRAINT "user_segment_fk1" FOREIGN KEY ("seg_id") REFERENCES "segment"("seg_id");