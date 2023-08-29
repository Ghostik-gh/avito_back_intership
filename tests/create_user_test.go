package tests

import (
	"avito_back_intership/internal/http-server/handlers/user/create_user"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

func TestCreateUser(t *testing.T) {

	testCases := []struct {
		name   string
		userID int
	}{
		{
			name:   "Create user without segment",
			userID: gofakeit.Year(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			e := httpexpect.Default(t, u.String())

			e.POST("/user/{user_id}").WithPath("user_id", tc.userID).WithBytes([]byte("{}")).Expect().Status(200)

			e.GET("/user/{user_id}").WithPath("user_id", tc.userID).Expect().Status(200).JSON().Object().HasValue("Segments", nil)

			e.DELETE("/user/{user_id}").WithPath("user_id", tc.userID).Expect().Status(200).JSON().Object().HasValue("status", "OK")
		})
	}
}

func TestCreateUserAndSegment(t *testing.T) {

	testCases := []struct {
		name       string
		segment    string
		procentage float64
		userID     int
		code       int
		error      string
	}{
		{
			name:    "Create user with segment",
			segment: gofakeit.Word(),
			userID:  gofakeit.Year(),
			code:    200,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			e := httpexpect.Default(t, u.String())

			e.POST("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("status", "OK")

			var input create_user.Request

			input.AddedSeg = append(input.AddedSeg, create_user.SegmentWithTime{Segment: tc.segment})

			inputByte, _ := json.Marshal(input)

			e.POST("/user/{user_id}").WithPath("user_id", tc.userID).WithBytes(inputByte).Expect().Status(200)

			e.GET("/user/{user_id}").WithPath("user_id", tc.userID).Expect().Status(200).JSON().Object().HasValue("Segments", []string{tc.segment})

			e.GET("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("userList", []string{fmt.Sprint(tc.userID)})

			e.DELETE("/user/{user_id}").WithPath("user_id", tc.userID).Expect().Status(200).JSON().Object().HasValue("status", "OK")

			e.DELETE("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("status", "OK")

		})
	}
}
