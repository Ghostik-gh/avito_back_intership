package create_user

import (
	"avito_back_intership/internal/lib/logger/sl"
	"database/sql"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	AddedSeg  []SegmentWithTime `json:"addedSeg,omitempty"`
	RemoveSeg []string          `json:"removeSeg,omitempty"`
}

type SegmentWithTime struct {
	Segment string `json:"segment,omitempty"`
	Time    string `json:"time,omitempty"`
}

type Response struct {
	response.Response
}

//go:generate mockery --name=URLSaver
type UserCreater interface {
	CreateUser(userID int) error
	CreateUserSegment(userID int, segment string) error
	DeleteUserSegment(userID int, segment string) error
	CreateLog(userID int, segName, opertaion string) error
	SegmentList() (*sql.Rows, error)
	UserInfo(userID int) (*sql.Rows, error)
	CreateUserSegmentTime(user_id int, segment, time string) error
}

// @Summary			Изменение сегментов у одного пользователя
// @Tags			User
// @Description		Изменяет состояние сегментов у пользователя, если пользователя нет, то он созадется. Также доступен timestamp в формате YYYY MM DD HH:MM:SS+00
// @ID				create-user
// @Accept			json
// @Produce			json
// @Param			user_id	path		int		true	"user id"
// @Param			input	body		Request	false	"user update data"
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/user/{user_id} [post]
func New(log *slog.Logger, userCreater UserCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.create_user.New"
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

		UserIDstr := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(UserIDstr)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, response.Error("user id not number"))
			return
		}

		if err := userCreater.CreateUser(userID); err != nil {
			log.Info("user already exists")
		} else {
			log.Info("user created")
		}

		userData, err := userCreater.UserInfo(userID)
		if err != nil {
			log.Error("failed to get user info", slog.String("user_id", UserIDstr))
			render.JSON(w, r, response.Error("failed to get user info"))
			return
		}
		defer userData.Close()

		var userSegmentList []string
		for userData.Next() {
			var row string
			userData.Scan(&row)
			userSegmentList = append(userSegmentList, row)
		}

		segments, err := userCreater.SegmentList()
		if err != nil {
			log.Error("failed to get segments info", slog.String("user_id", UserIDstr))
			render.JSON(w, r, response.Error("failed to get segments info"))
			return
		}
		defer segments.Close()

		var validSegments []string
		for segments.Next() {
			var row string
			segments.Scan(&row)
			if err != nil {
				log.Error(err.Error())
			}
			validSegments = append(validSegments, row)
		}

		for _, v := range req.AddedSeg {
			if slices.Contains(userSegmentList, v.Segment) {
				continue
			} else {
				log.Info("user already in this segment", slog.String("segment", v.Segment))
			}
			if slices.Contains(validSegments, v.Segment) {
				var err error
				if v.Time != "" {
					err = userCreater.CreateUserSegmentTime(userID, v.Segment, v.Time)
				} else {
					err = userCreater.CreateUserSegment(userID, v.Segment)
				}
				if err != nil {
					log.Error(err.Error())
					continue
				}
				err = userCreater.CreateLog(userID, v.Segment, "add")
				if err != nil {
					log.Error(err.Error())
				}
				userSegmentList = append(userSegmentList, v.Segment)
			} else {
				log.Info("segment not created", slog.String("segment", v.Segment))
			}
		}

		for _, v := range req.RemoveSeg {
			if slices.Contains(userSegmentList, v) {
				err := userCreater.DeleteUserSegment(userID, v)
				if err != nil {
					log.Error(err.Error())
				}
				userCreater.CreateLog(userID, v, "remove")
			}
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
