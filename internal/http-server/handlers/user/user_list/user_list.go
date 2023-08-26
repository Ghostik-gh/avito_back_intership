package user_list

import (
	"database/sql"
	"net/http"

	"log/slog"

	resp "avito_back_intership/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	UserList []string
	resp.Response
}

//go:generate mockery --name=URLSaver
type UserListGetter interface {
	UserList() (*sql.Rows, error)
}

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
			render.JSON(w, r, resp.Error("failed to get list of users"))
			return
		}

		var userList []string
		for rows.Next() {
			var one_row string
			rows.Scan(&one_row)
			userList = append(userList, one_row)
		}

		render.JSON(w, r, Response{
			UserList: userList,
			Response: resp.OK(),
		})

		// err = csv.CreateCSV(log, "users.csv", rows)
		// if err != nil {
		// 	log.Error("failed to create csv")
		// 	render.JSON(w, r, resp.Error("failed to create csv"))
		// 	return
		// }

	}

}
