package segment_list

import (
	"avito_back_intership/internal/lib/api/response"
	"database/sql"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	SegmentList []string
	response.Response
}

//go:generate mockery --name=URLSaver
type SegmentListGetter interface {
	SegmentList() (*sql.Rows, error)
}

// @Summary			Получения списка всех сегментов
// @Tags			Segment
// @Description		Возвращает список всех зарегистрированных пользователей
// @ID				segment-list
// @Accept			json
// @Produce			json
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/segment [get]
func New(log *slog.Logger, segmentListGetter SegmentListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.segment_list_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rows, err := segmentListGetter.SegmentList()
		if err != nil {
			log.Error("failed to get list of segments")
			render.JSON(w, r, response.Error("failed to get list of segments"))
			return
		}

		var segmentList []string
		for rows.Next() {
			var one_row string
			rows.Scan(&one_row)
			segmentList = append(segmentList, one_row)
		}

		render.JSON(w, r, Response{
			SegmentList: segmentList,
			Response:    response.OK(),
		})
	}
}
