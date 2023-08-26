package user_segments

import (
	"database/sql"
	"net/http"
	"strconv"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Segments []string
	resp.Response
}

//go:generate mockery --name=URLSaver
type UserGetter interface {
	UserInfo(int) (*sql.Rows, error)
}

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
			render.JSON(w, r, resp.Error("user id not number"))
			return
		}

		userData, err := userGetter.UserInfo(userID)
		if err != nil {
			log.Error("failed to get user info", slog.String("user_id", userStr))
			render.JSON(w, r, resp.Error("failed to get user info"))
			return
		}

		var userSegmentList []string
		for userData.Next() {
			var one_row string
			userData.Scan(&one_row)
			userSegmentList = append(userSegmentList, one_row)
		}

		render.JSON(w, r, Response{
			Segments: userSegmentList,
			Response: resp.OK(),
		})
	}

}
