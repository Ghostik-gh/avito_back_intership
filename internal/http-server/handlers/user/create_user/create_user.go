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
	CreateUser(userID int) error
	CreateUserSegment(userID int, segment string) error
	DeleteUserSegment(userID int, segment string) error
	CreateLog(userID int, seg_name, opertaion string) error
	SegmentList() (*sql.Rows, error)
	UserInfo(userID int) (*sql.Rows, error)
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

		userID, err := strconv.Atoi(req.UserID)
		if err != nil {
			log.Error("user id not number")
			render.JSON(w, r, resp.Error("user id not number"))
			return
		}

		if err := userCreater.CreateUser(userID); err != nil {
			log.Info("user already exists")
		} else {
			log.Info("user created")
		}

		userData, err := userCreater.UserInfo(userID)
		if err != nil {
			log.Error("failed to get user info", slog.String("user_id", req.UserID))
			render.JSON(w, r, resp.Error("failed to get user info"))
			return
		}

		var userSegmentList []string
		for userData.Next() {
			var row string
			userData.Scan(&row)
			userSegmentList = append(userSegmentList, row)
		}

		segments, err := userCreater.SegmentList()
		if err != nil {
			log.Error("failed to get segments info", slog.String("user_id", req.UserID))
			render.JSON(w, r, resp.Error("failed to get segments info"))
			return
		}
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
			if slices.Contains(userSegmentList, v) {
				continue
			}
			if slices.Contains(validSegments, v) {
				err := userCreater.CreateUserSegment(userID, v)
				if err != nil {
					log.Error(err.Error())
				}
				err = userCreater.CreateLog(userID, v, "add")
				if err != nil {
					log.Error(err.Error())
				}
				userSegmentList = append(userSegmentList, v)
			} else {
				log.Error("failed add segment to user", slog.String("segment", v))
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
			Response: resp.OK(),
		})
	}

}
