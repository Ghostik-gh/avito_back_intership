package create_segment

import (
	"avito_back_intership/internal/lib/logger/sl"
	"avito_back_intership/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

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

// @Success			200		{object}	userModel.TokenAccessModel	"data"
// @Failure			400,404	{object}	httpModel.ResponseMessage
// @Failure			500		{object}	httpModel.ResponseMessage
// @Failure			default	{object}	httpModel.ResponseMessage
//
//go:generate mockery --name=URLSaver
type SegmentCreator interface {
	CreateSegment(name string, amount string) error
	UserList() (*sql.Rows, error)
	CreateLog(user_id int, seg_name, opertaion string) error
	CreateUserSegment(user_id int, segment string) error
}

// @Summary			Создание сегмента
// @Tags			Segment
// @Description		Создание сегмента
// @ID				segment-creation
// @Accept			json
// @Produce			json
// @Param			input	body		Request						true	"segment name"
// @Router			/segment [post]
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

		amount, err := strconv.ParseFloat(req.Amount, 64)
		if err != nil {
			log.Info("amount not number ", sl.Err(err))
		}
		fmt.Printf("amount: %v\n", amount)
		if amount <= 0 || req.Amount == "" || amount > 100 {
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
		} else {
			err := segmentCreator.CreateSegment(req.Segment, req.Amount)
			if errors.Is(err, storage.ErrSegmentExists) {
				log.Error(storage.ErrSegmentExists.Error(), slog.String("segment", req.Segment))
				render.JSON(w, r, resp.Error(storage.ErrSegmentExists.Error()))
				return
			}
			rows, err := segmentCreator.UserList()
			if err != nil {
				log.Error(err.Error())
			}
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
			fmt.Printf("rand_users: %v\n", rand_users)
			for _, v := range rand_users {
				id, err := strconv.Atoi(v)
				if err != nil {
					log.Error(err.Error())
				}
				err = segmentCreator.CreateUserSegment(id, req.Segment)
				if err != nil {
					log.Error(err.Error())
				}
				err = segmentCreator.CreateLog(id, req.Segment, "add")
				if err != nil {
					log.Error(err.Error())
				}
			}
		}

		log.Info("segment created")
		render.JSON(w, r, Response{
			Response: resp.OK(),
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
