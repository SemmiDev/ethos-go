package errors_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/common/errors"
)

func TestCommonErrors(t *testing.T) {
	t.Parallel()

	Convey("Given the common errors package", t, func() {

		Convey("When creating IncorrectInputError", func() {
			err := errors.NewIncorrectInputError("invalid input", "test-slug")

			Convey("Then Error() should return the message", func() {
				So(err.Error(), ShouldEqual, "invalid input")
			})

			Convey("Then Slug() should return the slug", func() {
				So(err.Slug(), ShouldEqual, "test-slug")
			})
		})

		Convey("When creating NotFoundError", func() {
			err := errors.NewNotFoundError("resource not found", "not-found-slug")

			Convey("Then Error() should return the message", func() {
				So(err.Error(), ShouldEqual, "resource not found")
			})

			Convey("Then Slug() should return the slug", func() {
				So(err.Slug(), ShouldEqual, "not-found-slug")
			})
		})

		Convey("When creating UnauthorizedError", func() {
			err := errors.NewUnauthorizedError("access denied", "unauthorized-slug")

			Convey("Then Error() should return the message", func() {
				So(err.Error(), ShouldEqual, "access denied")
			})

			Convey("Then Slug() should return the slug", func() {
				So(err.Slug(), ShouldEqual, "unauthorized-slug")
			})
		})

		Convey("When creating ConflictError", func() {
			err := errors.NewConflictError("resource conflict", "conflict-slug")

			Convey("Then Error() should return the message", func() {
				So(err.Error(), ShouldEqual, "resource conflict")
			})

			Convey("Then Slug() should return the slug", func() {
				So(err.Slug(), ShouldEqual, "conflict-slug")
			})

			Convey("Then HTTPStatusCode should be 409", func() {
				So(err.HTTPStatusCode(), ShouldEqual, 409)
			})
		})
	})
}
