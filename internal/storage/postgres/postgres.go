package postgres

import (
	"avito_back_intership/internal/storage"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS
		"segment" (
			"name" VARCHAR(255) NOT NULL UNIQUE,
			"amount" FLOAT,
			CONSTRAINT "segments_pk" PRIMARY KEY ("name")
		)
		WITH (OIDS = FALSE);

		CREATE TABLE IF NOT EXISTS
			"people" (
				"user_id" integer NOT NULL,
				CONSTRAINT "user_pk" PRIMARY KEY ("user_id")
			)
		WITH (OIDS = FALSE);

		CREATE TABLE IF NOT EXISTS
			"user_segment" (
				"user_id" integer NOT NULL,
				"seg_name" VARCHAR(255) NOT NULL,
				"delete_time" TIMESTAMPTZ
			)
		WITH (OIDS = FALSE);

		CREATE TABLE IF NOT EXISTS
			"log" (
				"user_id" integer NOT NULL,
				"seg_name" VARCHAR(255) NOT NULL,
				"operation" VARCHAR(255) NOT NULL,
				"op_time" TIMESTAMPTZ NOT NULL
			)
		WITH (OIDS = FALSE);

		ALTER TABLE "user_segment" DROP CONSTRAINT IF EXISTS user_segment_fk0;
		ALTER TABLE "user_segment" DROP CONSTRAINT IF EXISTS user_segment_fk1;
		
		ALTER TABLE "user_segment"
		ADD
			CONSTRAINT "user_segment_fk0" FOREIGN KEY ("user_id") REFERENCES "people"("user_id") ON DELETE CASCADE;

		ALTER TABLE "user_segment"
		ADD
			CONSTRAINT "user_segment_fk1" FOREIGN KEY ("seg_name") REFERENCES "segment"("name") ON DELETE CASCADE;
			`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateSegment(name string, amount string) error {
	const op = "storage.postgres.CreateSegment"
	_, err := s.db.Exec(`INSERT INTO segment VALUES ($1, $2);`, name, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
	}
	return nil
}

func (s *Storage) DeleteSegment(name string) error {
	const op = "storage.postgres.DeleteSegment"

	res, err := s.db.Exec(`DELETE FROM segment WHERE name=$1;`, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("%s: %v", op, storage.ErrNothingDelete)
	}
	return nil
}

func (s *Storage) CreateUser(user_id int) error {
	const op = "storage.postgres.CreateUser"

	_, err := s.db.Exec(`INSERT INTO people VALUES ($1);`, user_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUser(user_id string) error {
	const op = "storage.postgres.DeleteUser"

	res, err := s.db.Exec(`DELETE FROM people WHERE user_id=$1;`, user_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("%s: %v", op, storage.ErrNothingDeleteUser)
	}

	return nil
}

func (s *Storage) CreateUserSegment(user_id int, segment string) error {
	const op = "storage.postgres.CreateUserSegment"

	_, err := s.db.Exec(`INSERT INTO user_segment (user_id, seg_name) VALUES ($1, $2);`, user_id, segment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateUserSegmentTime(user_id int, segment, time string) error {
	const op = "storage.postgres.CreateUserSegment"

	_, err := s.db.Exec(`INSERT INTO user_segment (user_id, seg_name, delete_time) VALUES ($1, $2, $3);`, user_id, segment, time)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUserSegment(user_id int, segment string) error {
	const op = "storage.postgres.DeleteUserSegment"
	_, err := s.db.Exec(`DELETE FROM user_segment WHERE user_id=$1 AND seg_name=$2;`, user_id, segment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateLog(user_id int, seg_name, opertaion string) error {
	const op = "storage.postgres.CreateLog"
	_, err := s.db.Exec(`INSERT INTO log VALUES ($1, $2, $3, now());`, user_id, seg_name, opertaion)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UserInfo(user_id int) (*sql.Rows, error) {
	const op = "storage.postgres.UserInfo"
	data, err := s.db.Query(`SELECT seg_name FROM user_segment WHERE user_id=$1 `, user_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (s *Storage) UserList() (*sql.Rows, error) {
	const op = "storage.postgres.UserList"
	data, err := s.db.Query(`SELECT user_id FROM people`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (s *Storage) UserLog(user_id int) (*sql.Rows, error) {
	const op = "storage.postgres.UserLog"
	data, err := s.db.Query(`SELECT * FROM log WHERE user_id=$1`, user_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (s *Storage) SegmentInfo(segment string) (*sql.Rows, error) {
	const op = "storage.postgres.SegmentInfo"
	data, err := s.db.Query(`SELECT user_id FROM user_segment WHERE seg_name=$1 `, segment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}
func (s *Storage) SegmentList() (*sql.Rows, error) {
	const op = "storage.postgres.SegmentList"
	data, err := s.db.Query(`SELECT name FROM segment`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

func (s *Storage) DeleteTTL() error {
	const op = "storage.postgres.DeleteTTL"
	_, err := s.db.Query(`DELETE FROM user_segment WHERE delete_time <= now()`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
