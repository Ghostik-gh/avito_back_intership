package create_user

import (
	"avito_back_intership/internal/lib/logger/sl"
	"database/sql"
	"net/http"
	"slices"
	"strconv"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	UserID    string   `json:"userID" validate:"required,number"`
	AddedSeg  []string `json:"addedSeg,omitempty"`
	RemoveSeg []string `json:"removeSeg,omitempty"`
}

type Response struct {
	resp.Response
}

//go:generate mockery --name=URLSaver
type UserCreater interface {
	CreateUser(user_id int) error
	CreateUserSegment(user_id int, segment string) error
	DeleteUserSegment(user_id int, segment string) error
	CreateLog(user_id int, seg_name, opertaion string) error
	SegmentList() (*sql.Rows, error)
	UserInfo(user_id int) (*sql.Rows, error)
}

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

		user_id, err := strconv.Atoi(req.UserID)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, resp.Error("user id not number"))
			return
		}

		if err := userCreater.CreateUser(user_id); err != nil {
			log.Info("user already exists")
			// log.Error("failed to create user", slog.String("user_id", req.UserID))
			// render.JSON(w, r, resp.Error("failed to create user"))
			// return
		} else {
			log.Info("user created")
		}

		user_data, err := userCreater.UserInfo(user_id)
		if err != nil {
			log.Error("failed to get user info", slog.String("user_id", req.UserID))
			render.JSON(w, r, resp.Error("failed to get user info"))
			return
		}

		var user_segment_list []string
		for user_data.Next() {
			var one_row string
			user_data.Scan(&one_row)
			user_segment_list = append(user_segment_list, one_row)
		}

		segments, err := userCreater.SegmentList()
		if err != nil {
			log.Error("failed to get segments info", slog.String("user_id", req.UserID))
			render.JSON(w, r, resp.Error("failed to get segments info"))
			return
		}
		var valid_segments []string
		for segments.Next() {
			var one_row string
			segments.Scan(&one_row)
			if err != nil {
				log.Error(err.Error())
			}
			valid_segments = append(valid_segments, one_row)
		}

		for _, v := range req.AddedSeg {
			if slices.Contains(user_segment_list, v) {
				continue
			}
			if slices.Contains(valid_segments, v) {
				err := userCreater.CreateUserSegment(user_id, v)
				if err != nil {
					log.Error(err.Error())
				}
				err = userCreater.CreateLog(user_id, v, "add")
				if err != nil {
					log.Error(err.Error())
				}
				user_segment_list = append(user_segment_list, v)
			} else {
				log.Error("failed add segment to user", slog.String("segment", v))
			}
		}

		for _, v := range req.RemoveSeg {
			if slices.Contains(user_segment_list, v) {
				err := userCreater.DeleteUserSegment(user_id, v)
				if err != nil {
					log.Error(err.Error())
				}
				userCreater.CreateLog(user_id, v, "remove")
			}
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}

}
