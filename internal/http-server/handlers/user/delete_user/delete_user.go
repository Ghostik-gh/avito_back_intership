package delete_user

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
	UserID string `json:"user_id" validate:"required,number"`
}

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type UserDeleter interface {
	DeleteUser(user_id string) error
}

func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.deleter.New"
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

		if err = userDeleter.DeleteUser(req.UserID); err != nil {
			log.Error("failed to delete user", slog.String("user_id", req.UserID))
			render.JSON(w, r, resp.Error("failed to delete user"))
			return
		}
		log.Info("user deleted")
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
