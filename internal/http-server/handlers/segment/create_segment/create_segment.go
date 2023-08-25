package create_segment

import (
	"avito_back_intership/internal/lib/logger/sl"
	"avito_back_intership/internal/storage"
	"errors"
	"net/http"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Segment string `json:"segment" validate:"required"`
	Amount  string `json:"amount,omitempty"`
}

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type SegmentCreator interface {
	CreateSegment(name string, amount string) error
}

func New(log *slog.Logger, segmentCreator SegmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.create_segment.New"
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

		if err := segmentCreator.CreateSegment(req.Segment, "0"); err != nil {
			if errors.Is(err, storage.ErrSegmentExists) {
				log.Error(storage.ErrSegmentExists.Error(), slog.String("segment", req.Segment))
				render.JSON(w, r, resp.Error(storage.ErrSegmentExists.Error()))
				return
			}
			log.Error("failed to create segment", slog.String("segment", req.Segment))
			render.JSON(w, r, resp.Error("failed to create segment"))
			return
		}
		log.Info("segment created")

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
