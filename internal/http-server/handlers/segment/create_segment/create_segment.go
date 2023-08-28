package create_segment

import (
	"avito_back_intership/internal/storage"
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"strconv"

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
type SegmentCreator interface {
	CreateSegment(name string, amount string) error
	UserList() (*sql.Rows, error)
	CreateLog(user_id int, seg_name, opertaion string) error
	CreateUserSegment(user_id int, segment string) error
}

// @Summary			Создание сегмента
// @Tags			Segment
// @Description		Создает сегмент, необязательное поле указывает на то сколько процентов (float от 0 до 100) от зарегистрированных пользователей попадет в выборка
// @ID				segment-creation
// @Accept			json
// @Produce			json
// @Param			segment		path		string	true	"segment name"
// @Param			percentage	path		number	false	"percentage"
// @Success			200			{object}	Response
// @Failure			default		{object}	Response
// @Router			/segment/{segment}/{percentage} [post]
func New(log *slog.Logger, segmentCreator SegmentCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.create_segment.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		segment := chi.URLParam(r, "segment")
		percentage := chi.URLParam(r, "percentage")
		var err error
		var amount float64
		amount, err = strconv.ParseFloat(percentage, 64)
		if amount == 0 {
			err = nil
			percentage = "0"
		}
		if err != nil {
			log.Error("not number", slog.String("percentage", percentage))
			render.JSON(w, r, response.Error("not number"))
			return
		} else {
			if amount < 0 || amount > 100 {
				log.Error("wrong number", slog.String("segment", segment))
				render.JSON(w, r, response.Error("wrong number"))
				return
			}

			err := segmentCreator.CreateSegment(segment, percentage)
			if errors.Is(err, storage.ErrSegmentExists) {
				log.Error(storage.ErrSegmentExists.Error(), slog.String("segment", segment))
				render.JSON(w, r, response.Error(storage.ErrSegmentExists.Error()))
				return
			}
			rows, err := segmentCreator.UserList()
			if err != nil {
				log.Error(err.Error())
			}
			defer rows.Close()
			var users []string
			for rows.Next() {
				var row string
				rows.Scan(&row)
				if err != nil {
					log.Error(err.Error())
				}
				users = append(users, row)
			}
			var count int = int(float64(len(users)) * (amount / 100))
			rand_users := chooseRandomUsers(users, count)

			for _, v := range rand_users {
				id, err := strconv.Atoi(v)
				if err != nil {
					log.Error(err.Error())
				}
				err = segmentCreator.CreateUserSegment(id, segment)
				if err != nil {
					log.Error(err.Error())
				}
				err = segmentCreator.CreateLog(id, segment, "add")
				if err != nil {
					log.Error(err.Error())
				}
			}
		}

		log.Info("segment created")
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}

func chooseRandomUsers(userIDs []string, count int) []string {
	if count > len(userIDs) {
		count = len(userIDs)
	}
	chosenUsers := make(map[string]bool)

	for len(chosenUsers) < count {
		randomIndex := rand.Intn(len(userIDs))
		chosenUsers[userIDs[randomIndex]] = true
	}
	result := make([]string, 0, count)
	for user := range chosenUsers {
		result = append(result, user)
	}

	return result
}
