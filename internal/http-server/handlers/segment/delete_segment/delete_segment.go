package delete_segment

import (
	"avito_back_intership/internal/storage"
	"errors"
	"fmt"
	"net/http"

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
type SegmentDeleter interface {
	DeleteSegment(name string) error
}

// @Summary			Удаление сегмента
// @Tags			Segment
// @Description		Удаление сегмента
// @ID				segment-deletion
// @Accept			json
// @Produce			json
// @Param			name	path		string						true	"segment name"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/segment/{segment} [delete]
func New(log *slog.Logger, segmentDeleter SegmentDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.segment_deleter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		segment := chi.URLParam(r, "segment")

		fmt.Printf("segment: %v\n", segment)
		if err := segmentDeleter.DeleteSegment(segment); err != nil {
			if errors.Is(err, storage.ErrNothingDelete) {
				log.Error(storage.ErrNothingDelete.Error(), slog.String("segment", segment))
				render.JSON(w, r, response.Error(storage.ErrNothingDelete.Error()))
				return
			}
			log.Error("failed to delete segment", slog.String("segment", segment))
			render.JSON(w, r, response.Error("failed to delete segment"))
			return
		}

		log.Info("segment deleted")
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}

}
