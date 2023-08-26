package delete_user

import (
	"avito_back_intership/internal/storage"
	"net/http"

	"log/slog"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	UserID string `json:"user_id" validate:"required,number"`
}

type Response struct {
	response.Response
}

//go:generate mockery --name=URLSaver
type UserDeleter interface {
	DeleteUser(user_id string) error
}

// @Summary			Удаление пользователя
// @Tags			User
// @Description		Удаление пользователя, удаляются все записи
// @ID				user-deletion
// @Accept			json
// @Produce			json
// @Param			user_id	path		int			true	"user id"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/user/{user_id} [delete]
func New(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.deleter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		req.UserID = chi.URLParam(r, "user_id")

		if err := userDeleter.DeleteUser(req.UserID); err != nil {
			log.Error(storage.ErrNothingDeleteUser.Error(), slog.String("user", req.UserID))
			render.JSON(w, r, response.Error(storage.ErrNothingDeleteUser.Error()))
			return
		}

		log.Info("user deleted")
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
