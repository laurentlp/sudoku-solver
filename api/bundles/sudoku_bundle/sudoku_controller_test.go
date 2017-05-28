package sudokubundle_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/laurentlp/sudoku-solver/api/bundles/sudoku_bundle"
	"github.com/laurentlp/sudoku-solver/api/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestControllerSpec(t *testing.T) {
	Convey("Given a running server and a controller instance", t, func() {
		// load error messages
		if err := errors.LoadErrors("../../errors/error_templates.yaml"); err != nil {
			panic(fmt.Errorf("Failed to read the error message file: %s", err))
		}

		c := sudokubundle.SudokuController{}

		// Create and start the sudoku
		mux := http.NewServeMux()
		server := httptest.NewServer(mux)
		defer server.Close()

		Convey("When Solve is called from handler with a good sudoku", func() {
			mux.HandleFunc("/sudoku", c.Solve)

			reader := strings.NewReader(`
				{
					"sudoku" : "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
				}
			`)

			resp, err := http.Post(server.URL+"/sudoku", "application/json", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 200 with correct JSON response", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
				So(string(body), ShouldEqual, `{"sudoku":"417369825632158947958724316825437169791586432346912758289643571573291684164875293","solved":true}`)
			})
		})

		Convey("When Solve is called from handler with an short sudoku", func() {
			mux.HandleFunc("/sudoku", c.Solve)

			reader := strings.NewReader(`{"sudoku": "4.....8.5.3..........7......2.....6.....8.4......1......6.3.7.5..2.....1.4......"}`)

			resp, err := http.Post(server.URL+"/sudoku", "", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400 with correct JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(string(body), ShouldEqual, `{"error_code":"BAD_REQUEST","message":"Invalid grid size: expected grid size of 81 found grid size of 80"}`)
			})
		})

		Convey("When Solve is called from handler with nothing", func() {
			mux.HandleFunc("/sudoku", c.Solve)

			reader := strings.NewReader(``)

			resp, err := http.Post(server.URL+"/sudoku", "", reader)
			if err != nil {
				t.Fatal(err)
			}

			Convey("Then response should be 400 with correct JSON error", func() {
				body, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
				So(string(body), ShouldEqual, `{"error_code":"BAD_REQUEST","message":"No information was sent to the server. Please send a valid sudoku."}`)
			})
		})
	})
}
