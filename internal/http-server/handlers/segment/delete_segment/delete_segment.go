package delete_segment

import (
	"avito_back_intership/internal/lib/logger/sl"
	"net/http"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Segment string `json:"segment" validate:"required"`
}

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type SegmentDeleter interface {
	DeleteSegment(name string) (int64, error)
}

func New(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.segment_deleter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body request", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		var rows int64
		if rows, err = segmentDeleter.DeleteSegment(req.Segment); err != nil {
			log.Info("failed to delete segment", slog.String("segment", req.Segment))
			render.JSON(w, r, resp.Error("failed to delete segment"))
			return
		}
		if rows == 0 {
			log.Info("nothing to delete")
			render.JSON(w, r, Response{
				Response: resp.Error("nothing to delete"),
			})
		} else {
			log.Info("segment deleted")
			render.JSON(w, r, Response{
				Response: resp.OK(),
			})
		}
	}

}
