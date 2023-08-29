package user_log

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
}

//go:generate mockery --name=URLSaver
type UserLogGetter interface {
	UserLog(user_id int) (*sql.Rows, error)
}

// @Summary			Лог пользователя
// @Tags			Log
// @Description		Возвращает CSV файл для выбранного пользователя
// @ID				user-log
// @Produce			octet-stream
// @Param			user_id	path		int		true	"user id"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/log/{user_id} [get]
func New(log *slog.Logger, userLogGetter UserLogGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user_log.getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userStr := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(userStr)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, response.Error("user id not number"))
			return
		}

		rows, err := userLogGetter.UserLog(userID)
		if err != nil {
			log.Error("failed to get user's log")
			render.JSON(w, r, response.Error("failed to get user's log"))
			return
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprint("attachment;filename="+userStr+".csv"))

		writer := csv.NewWriter(w)

		columns, err := rows.Columns()
		if err != nil {
			log.Error(err.Error())
		}
		writer.Write(columns)

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
		rows.Close()

		writer.Flush()
		w.WriteHeader(http.StatusOK)
		log.Info("csv file created")
	}
}
