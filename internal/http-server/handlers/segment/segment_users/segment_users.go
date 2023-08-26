package segment_users

import (
	"database/sql"
	"net/http"

	"log/slog"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// type Request struct {
// 	Segment string `json:"segment" validate:"required"`
// }

type Response struct {
	UserList []string `json:"userList"`
	response.Response
}

//go:generate mockery --name=URLSaver
type SegmentGetter interface {
	SegmentInfo(segment string) (*sql.Rows, error)
}

// @Summary			Получение всех пользователей в данном сегменте
// @Tags			Segment
// @Description		Получение всех пользователей в данном сегменте
// @ID				segment-user-list
// @Accept			json
// @Produce			json
// @Param			segment	path		string						true	"segment name"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/segment/{segment} [get]
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
			render.JSON(w, r, response.Error("empty request"))
			return
		}

		rows, err := segmentGetter.SegmentInfo(segment)
		if err != nil {
			log.Error("failed to get users in segment", slog.String("segment", segment))
			render.JSON(w, r, response.Error("failed to get users in segment "+segment))
			return
		}

		var userList []string
		for rows.Next() {
			var tmp string
			rows.Scan(&tmp)
			userList = append(userList, tmp)
		}

		render.JSON(w, r, Response{
			UserList: userList,
			Response: response.OK(),
		})

		// err = csv.CreateCSV(log, "segments_for_user.csv", rows)
		// if err != nil {
		// 	log.Error("failed to create csv", slog.String("segment", segment))
		// 	render.JSON(w, r, response.Error("failed to create csv "+segment))
		// 	return
		// }

		// log.Info("csv file created")
		// render.JSON(w, r, Response{
		// 	Response: response.OK(),
		// })
	}

}
