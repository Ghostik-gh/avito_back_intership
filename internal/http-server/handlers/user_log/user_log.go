package user_log

import (
	"database/sql"
	"net/http"
	"strconv"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"
	"avito_back_intership/internal/lib/csv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type UserLogGetter interface {
	UserLog(user_id int) (*sql.Rows, error)
}

func New(log *slog.Logger, userLogGetter UserLogGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.user_log_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		user_str := chi.URLParam(r, "user_id")

		user_id, err := strconv.Atoi(user_str)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, resp.Error("user id not number"))
			return
		}

		rows, err := userLogGetter.UserLog(user_id)
		if err != nil {
			log.Error("failed to get user's log")
			render.JSON(w, r, resp.Error("failed to get user's log"))
			return
		}

		err = csv.CreateCSV(log, "user_log.csv", rows)
		if err != nil {
			log.Error("failed to create csv")
			render.JSON(w, r, resp.Error("failed to create csv"))
			return
		}

		log.Info("csv file created")
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
