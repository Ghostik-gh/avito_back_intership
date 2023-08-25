package user_list

import (
	"database/sql"
	"net/http"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"
	"avito_back_intership/internal/lib/csv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type UsertListGetter interface {
	UserList() (*sql.Rows, error)
}

func New(log *slog.Logger, usertListGetter UsertListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.usert_list_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rows, err := usertListGetter.UserList()
		if err != nil {
			log.Error("failed to get list of users")
			render.JSON(w, r, resp.Error("failed to get list of users"))
			return
		}

		err = csv.CreateCSV(log, "1.csv", rows)
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
