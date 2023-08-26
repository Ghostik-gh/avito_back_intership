package delete_segment

import (
	"avito_back_intership/internal/lib/logger/sl"
	"avito_back_intership/internal/storage"
	"errors"
	"net/http"

	"log/slog"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Segment string `json:"segment" validate:"required"`
}

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
// @Param			input	body		Request						true	"segment name"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/segment [delete]
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
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}
		if err = segmentDeleter.DeleteSegment(req.Segment); err != nil {
			if errors.Is(err, storage.ErrNothingDelete) {
				log.Error(storage.ErrNothingDelete.Error(), slog.String("segment", req.Segment))
				render.JSON(w, r, response.Error(storage.ErrNothingDelete.Error()))
				return
			}

			log.Error("failed to delete segment", slog.String("segment", req.Segment))
			render.JSON(w, r, response.Error("failed to delete segment"))
			return
		}
		log.Info("segment deleted")
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}

}
