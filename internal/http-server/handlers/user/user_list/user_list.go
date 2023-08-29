package user_list

import (
	"database/sql"
	"net/http"

	"log/slog"

	"avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	UserList []string
	response.Response
}

//go:generate mockery --name=URLSaver
type UserListGetter interface {
	UserList() (*sql.Rows, error)
}

// @Summary			Список всех пользователей
// @Tags			User
// @Description		Возвращает список всех зарегистрированных пользователей
// @ID				user-list
// @Accept			json
// @Produce			json
// @Success			200		{object}	Response
// @Failure			default	{object}	Response
// @Router			/user [get]
func New(log *slog.Logger, userListGetter UserListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.list_getter.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rows, err := userListGetter.UserList()
		if err != nil {
			log.Error("failed to get list of users")
			render.JSON(w, r, response.Error("failed to get list of users"))
			return
		}

		var userList []string
		for rows.Next() {
			var one_row string
			rows.Scan(&one_row)
			userList = append(userList, one_row)
		}
		rows.Close()

		render.JSON(w, r, Response{
			UserList: userList,
			Response: response.OK(),
		})
	}
}
