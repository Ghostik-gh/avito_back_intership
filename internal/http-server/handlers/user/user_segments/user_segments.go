package user_segments

import (
	"database/sql"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"avito_back_intership/internal/lib/api/response"
	"avito_back_intership/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Segments []string
	response.Response
}

//go:generate mockery --name=URLSaver
type UserGetter interface {
	UserInfo(int) (*sql.Rows, error)
	UserList() (*sql.Rows, error)
}

// @Summary			Сегменты пользователя
// @Tags			User
// @Description		Получение всех сегментов данного пользователя
// @ID				user-segment-list
// @Accept			json
// @Produce			json
// @Param			user_id	path		int	true	"user id"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/user/{user_id} [get]
func New(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.user_segments.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userStr := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(userStr)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, response.Error("user id not number"))
			return
		}

		rows, err := userGetter.UserList()
		if err != nil {
			log.Error("failed to get user list", slog.String("user_id", userStr))
			render.JSON(w, r, response.Error("failed to get user list"))
			return
		}
		var userList []string
		for rows.Next() {
			var row string
			rows.Scan(&row)
			if err != nil {
				log.Error(err.Error())
			}
			userList = append(userList, row)
		}
		rows.Close()

		if !slices.Contains(userList, userStr) {
			log.Error(storage.ErrNothingDeleteUser.Error(), slog.String("user_id", userStr))
			render.JSON(w, r, response.Error(storage.ErrNothingDeleteUser.Error()))
			return
		}

		userData, err := userGetter.UserInfo(userID)
		if err != nil {
			log.Error("failed to get user info", slog.String("user_id", userStr))
			render.JSON(w, r, response.Error("failed to get user info"))
			return
		}
		defer userData.Close()

		var userSegmentList []string
		for userData.Next() {
			var one_row string
			userData.Scan(&one_row)
			userSegmentList = append(userSegmentList, one_row)
		}

		render.JSON(w, r, Response{
			Segments: userSegmentList,
			Response: response.OK(),
		})
	}
}
