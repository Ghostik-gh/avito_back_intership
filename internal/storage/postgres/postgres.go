package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sql.Open("postgres", "host=localhost port=5432 user=root password=qweasd sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
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
			`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateSegment(name string, amount float64) error {
	const op = "storage.postgres.CreateSegment"

	_, err := s.db.Exec(`INSERT INTO segment (name, amount) VALUES ('?', ?);`, name, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteSegment(name string) error {
	const op = "storage.postgres.DeleteSegment"

	_, err := s.db.Exec(`DELETE FROM segment WHERE name='?';`, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateUser(user_id int) error {
	const op = "storage.postgres.CreateUser"

	_, err := s.db.Exec(`INSERT INTO people VALUES (?);`, user_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUser(user_id int) error {
	const op = "storage.postgres.DeleteUser"

	_, err := s.db.Exec(`DELETE FROM people WHERE user_id=?;`, user_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateUserSegment(user_id int, segment string) error {
	const op = "storage.postgres.CreateUserSegment"

	_, err := s.db.Exec(`INSERT INTO user_segment (user_id, seg_name) VALUES (?, '?');`, user_id, segment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUserSegment(user_id int, segment string) error {
	const op = "storage.postgres.DeleteUserSegment"
	_, err := s.db.Exec(`DELETE FROM user_segment WHERE user_id=? AND seg_name='?';`, user_id, segment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateLog(user_id int, seg_name, opertaion string) error {
	const op = "storage.postgres.CreateLog"
	_, err := s.db.Exec(`INSERT INTO user_segment VALUES (?, '?', '?', now());`, user_id, seg_name, opertaion)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
