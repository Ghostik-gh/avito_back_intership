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
			CONSTRAINT "user_segment_fk1" FOREIGN KEY ("seg_id") REFERENCES "segment"("seg_id");`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateSegment(name string, amount float64) error {
	const op = "storage.postgres.CreateSegment"
	_, err := s.db.Exec(`
	INSERT INTO segment (name) VALUES ('?');`, name)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// func (s *Storage) SaveURL(urlLong string, alias string) error {
// 	const op = "storage.postgres.SaveURL"

// 	query, err := s.db.Prepare(context.Background(), "save url", "INSERT INTO url(url, alias) VALUES(?, ?)")
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	fmt.Printf("query.Fields: %v\n", query.Fields)
// 	// _, err = query.Exec(urlLong, alias)
// 	// if err != nil {
// 	// 	return fmt.Errorf("%s: %w", op, err)
// 	// }
// 	return nil

// }

// func (s *Storage) GetURL(alias string) (string, error) {
// 	const op = "storage.postgres.GetURL"

// 	query, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
// 	if err != nil {
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}
// 	var url string

// 	err = query.QueryRow(alias).Scan(&url)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return "", storage.ErrURLNotFound
// 		}
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}

// 	return url, nil
// }

// func (s *Storage) DeleteURL(alias string) error {
// 	const op = "storage.postgres.DeleteURL"

// 	query, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	result, err := query.Exec(alias)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrURLNotFound) {
// 			return storage.ErrURLNotFound
// 		}
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	if rows, _ := result.RowsAffected(); rows == 0 {
// 		return fmt.Errorf("%s: %v", op, "nothing to delete")
// 	}

// 	return nil
// }
