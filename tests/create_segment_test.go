package tests

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

func TestCreateSegment(t *testing.T) {

	testCases := []struct {
		name       string
		segment    string
		procentage float64
		code       int
		error      string
	}{
		{
			name:    "Segment",
			segment: gofakeit.Word(),
			code:    200,
		},
		{
			name:       "Segment With Procentage",
			segment:    gofakeit.Word(),
			procentage: gofakeit.Float64Range(0, 100),
			code:       200,
		},
		{
			name:       "Segment Wrong Procentage",
			segment:    gofakeit.Word(),
			procentage: gofakeit.Float64Range(101, 1000),
			code:       200,
			error:      "wrong number",
		},
		{
			name:       "Segment Wrong Procentage",
			segment:    gofakeit.Word(),
			procentage: -gofakeit.Float64Range(1, 1000),
			code:       200,
			error:      "wrong number",
		},
		{
			name:    "Segment Empty",
			segment: "",
			code:    404,
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
			var req *httpexpect.Request
			if tc.procentage != 0 {
				req = e.POST("/segment/{segment}/{procentage}").
					WithPath("segment", tc.segment).WithPath("procentage", tc.procentage)
			} else {
				req = e.POST("/segment/{segment}").WithPath("segment", tc.segment)
			}

			res := req.Expect().Status(tc.code)

			if tc.code != 200 {
				return
			}

			obj := res.JSON().Object()

			if tc.error != "" {
				obj.HasValue("error", tc.error)
				return
			} else {
				obj.HasValue("status", "OK")
			}

			e.POST("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("error", "segment exists")
			if tc.procentage == 0 {
				e.GET("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("userList", nil)
			}
			e.DELETE("/segment/{segment}").WithPath("segment", tc.segment).Expect().Status(200).JSON().Object().HasValue("status", "OK")
		})
	}
}
