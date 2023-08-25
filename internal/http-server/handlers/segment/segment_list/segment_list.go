package segment_list

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
type SegmentListGetter interface {
	SegmentList() (*sql.Rows, error)
}

func New(log *slog.Logger, segmentListGetter SegmentListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.segment_list_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rows, err := segmentListGetter.SegmentList()
		if err != nil {
			log.Error("failed to get list of segments")
			render.JSON(w, r, resp.Error("failed to get list of segments"))
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
