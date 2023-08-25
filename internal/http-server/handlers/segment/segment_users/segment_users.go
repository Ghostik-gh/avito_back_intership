package segment_users

import (
	"database/sql"
	"net/http"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"
	"avito_back_intership/internal/lib/csv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// type Request struct {
// 	Segment string `json:"segment" validate:"required"`
// }

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type SegmentGetter interface {
	SegmentInfo(segment string) (*sql.Rows, error)
}

func New(log *slog.Logger, segmentGetter SegmentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.segment_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		segment := chi.URLParam(r, "segment")

		if segment == "" {
			log.Error("empty request")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		rows, err := segmentGetter.SegmentInfo(segment)
		if err != nil {
			log.Error("failed to get users in segment", slog.String("segment", segment))
			render.JSON(w, r, resp.Error("failed to get users in segment "+segment))
			return
		}

		err = csv.CreateCSV(log, "1.csv", rows)
		if err != nil {
			log.Error("failed to create csv", slog.String("segment", segment))
			render.JSON(w, r, resp.Error("failed to create csv "+segment))
			return
		}

		log.Info("csv file created")
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
