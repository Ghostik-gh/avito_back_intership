package csv

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	"log/slog"
)

func CreateCSV(log *slog.Logger, name string, rows *sql.Rows) error {
	const op = "lib.csv.CreateCSV"
	log = log.With(
		slog.String("op", op),
	)
	// Создание файла CSV
	file, err := os.Create(name)
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()

	// Создание записывающего объекта CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Получение и запись заголовков столбцов
	columns, err := rows.Columns()
	if err != nil {
		log.Error(err.Error())
	}
	writer.Write(columns)

	// Чтение данных из *sql.Rows и их запись в файл CSV
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range columns {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			log.Error(err.Error())
		}
		row := make([]string, len(columns))
		for i, col := range values {
			row[i] = fmt.Sprintf("%v", col)
		}
		writer.Write(row)
	}
	if err := rows.Err(); err != nil {
		log.Error(err.Error())
	}

	return nil
}
