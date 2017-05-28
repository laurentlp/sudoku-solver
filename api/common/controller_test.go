package common_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/laurentlp/sudoku-solver/api/bundles/sudoku_bundle"
	"github.com/laurentlp/sudoku-solver/api/common"
	"github.com/laurentlp/sudoku-solver/api/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCommonController(t *testing.T) {
	Convey("Given a running server and a controller instance", t, func() {
		c := common.Controller{}

		// load error messages
		if err := errors.LoadErrors("../errors/error_templates.yaml"); err != nil {
			panic(fmt.Errorf("Failed to read the error message file: %s", err))
		}

		// Create and start the server
		mux := http.NewServeMux()
		server := httptest.NewServer(mux)
		defer server.Close()

		Convey("When SendJSON is called from handler with a code 200", func() {
			mux.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) {
				info := struct {
					Something string `json:"something"`
				}{
					Something: "with a message",
				}
				c.SendJSON(w, nil, &info, http.StatusOK)
			})

			resp, err := http.Get(server.URL + "/test1")
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 200 with the correct JSON", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(string(body), ShouldEqual, `{"something":"with a message"}`)
			})
		})

		Convey("When SendJSON is called from handler with invalid informations", func() {
			mux.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) {
				info := make(chan int)
				c.SendJSON(w, nil, &info, http.StatusOK)
			})

			resp, err := http.Get(server.URL + "/test2")
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 500 with the JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
				So(string(body), ShouldEqual, `{"error": "Internal server error"}`)
			})
		})

		Convey("When SendJSON is called from handler with nil value", func() {
			mux.HandleFunc("/test3", func(w http.ResponseWriter, r *http.Request) {
				c.SendJSON(w, nil, nil, http.StatusOK)
			})

			resp, err := http.Get(server.URL + "/test3")
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 500 with the JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
				So(string(body), ShouldEqual, `{"error": "Internal server error"}`)
			})
		})

		Convey("When SendJSON is called from handler with a code 400", func() {
			mux.HandleFunc("/test4", func(w http.ResponseWriter, r *http.Request) {
				info := struct {
					Something string `json:"something"`
				}{
					Something: "with a message",
				}
				c.SendJSON(w, nil, &info, http.StatusBadRequest)
			})

			resp, err := http.Get(server.URL + "/test4")
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400", func() {
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When MapJSON is called from handler with nothing", func() {
			mux.HandleFunc("/test5", func(w http.ResponseWriter, r *http.Request) {
				var model sudokubundle.Sudoku
				err := c.MapJSON(w, r, &model)
				c.SendJSON(w, nil, err, err.Status)
			})

			reader := strings.NewReader(``)

			resp, err := http.Post(server.URL+"/test5", "application/json", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400 with a JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(string(body), ShouldEqual, `{"error_code":"BAD_REQUEST","message":"No information was sent to the server. Please send a valid sudoku."}`)
			})
		})

		Convey("When MapJSON is called from handler with invalid JSON", func() {
			mux.HandleFunc("/test6", func(w http.ResponseWriter, r *http.Request) {
				var model sudokubundle.Sudoku
				err := c.MapJSON(w, r, &model)
				c.SendJSON(w, nil, err, err.Status)
			})

			reader := strings.NewReader(`{"aaaaa : aaaaa}`)

			resp, err := http.Post(server.URL+"/test6", "application/json", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400 with a JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(string(body), ShouldEqual, `{"error_code":"BAD_REQUEST","message":"The informations sent to the server contains errors."}`)
			})
		})

		Convey("When MapJSON is called from handler with to much data", func() {
			mux.HandleFunc("/test7", func(w http.ResponseWriter, r *http.Request) {
				var model sudokubundle.Sudoku
				err := c.MapJSON(w, r, &model)
				c.SendJSON(w, nil, err, err.Status)
			})

			reader := strings.NewReader(`{
				"sudoku": "624578139135496827789123456216385794857964213493217685942651378568732941371849562",
				"sudoku": "624578139135496827789123456216385794857964213493217685942651378568732941371849562",
				"sudoku": "624578139135496827789123456216385794857964213493217685942651378568732941371849562",
				"sudoku": "624578139135496827789123456216385794857964213493217685942651378568732941371849562",
				"solved": true}`)

			resp, err := http.Post(server.URL+"/test7", "application/json", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400 with a JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(string(body), ShouldEqual, `{"error_code":"BAD_REQUEST","message":"The informations sent to the server contains errors."}`)
			})
		})

		Convey("When MapJSON is called from handler with valid JSON", func() {
			mux.HandleFunc("/test8", func(w http.ResponseWriter, r *http.Request) {
				var model sudokubundle.Sudoku
				c.MapJSON(w, r, &model)
				c.SendJSON(w, nil, model, http.StatusOK)
			})

			json := `{"sudoku":"...57..3.1......2.7...234......8...4..7..4...49....6.5.42...3.....7..9....18.....","solved":false}`
			reader := strings.NewReader(json)

			resp, err := http.Post(server.URL+"/test8", "application/json", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 200 with the JSON response", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(string(body), ShouldEqual, json)
			})
		})
	})
}
