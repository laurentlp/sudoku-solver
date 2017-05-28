package errors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/laurentlp/sudoku-solver/api/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorHandler(t *testing.T) {
	Convey("Given an error template and handler", t, func() {
		Convey("When errors.LoadErrors is called from handler with an invalid path", func() {
			// load error messages
			err := errors.LoadErrors("../../errors/error_templates.yaml")

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})

		// load error messages
		if err := errors.LoadErrors("../errors/error_templates.yaml"); err != nil {
			panic(fmt.Errorf("Failed to read the error message file: %s", err))
		}

		Convey("When errors.BadRequest is called from handler with an error message", func() {
			msg := "A bad error occurred"
			err := errors.BadRequest(msg)

			Convey("Then error should have an HTTP status of 400 with the correct error message", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, msg)
				So(err.StatusCode(), ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When errors.BadRequest is called from handler with an empty message", func() {
			msg := ""
			err := errors.BadRequest(msg)

			Convey("Then error should have an HTTP status of 400 with the rmpty error message", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, msg)
				So(err.StatusCode(), ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}
