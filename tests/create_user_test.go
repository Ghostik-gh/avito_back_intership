package tests

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

func TestCreateUser(t *testing.T) {

	testCases := []struct {
		name       string
		segment    string
		procentage float64
		userID     int
		code       int
		error      string
	}{
		{
			name:   "Create user without segment",
			userID: gofakeit.Year(),
			code:   200,
		},
		// {
		// 	name:       "Segment With Procentage",
		// 	segment:    gofakeit.Word(),
		// 	procentage: gofakeit.Float64Range(0, 100),
		// 	code:       200,
		// },
		// {
		// 	name:       "Segment Wrong Procentage",
		// 	segment:    gofakeit.Word(),
		// 	procentage: gofakeit.Float64Range(101, 1000),
		// 	code:       200,
		// 	error:      "wrong number",
		// },
		// {
		// 	name:    "Segment Empty",
		// 	segment: "",
		// 	code:    404,
		// },
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			e := httpexpect.Default(t, u.String())

			e.POST("/user/{user_id}").WithPath("user_id", tc.userID).WithBytes([]byte("{}")).Expect().Status(tc.code)

			e.GET("/user/{user_id}").WithPath("user_id", tc.userID).Expect().Status(200).JSON().Object().HasValue("Segments", nil)

		})
	}
}
