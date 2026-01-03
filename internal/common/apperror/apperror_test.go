package apperror_test

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/common/apperror"
)

func TestAppError(t *testing.T) {
	t.Parallel()

	Convey("Given the apperror package", t, func() {

		Convey("When creating a new AppError", func() {
			underlyingErr := errors.New("underlying error")
			err := apperror.New("TEST_CODE", "test message", http.StatusBadRequest, underlyingErr)

			Convey("Then it should have the correct code", func() {
				So(err.Code, ShouldEqual, "TEST_CODE")
			})

			Convey("Then it should have the correct message", func() {
				So(err.Message, ShouldEqual, "test message")
			})

			Convey("Then it should have the correct status code", func() {
				So(err.StatusCode, ShouldEqual, http.StatusBadRequest)
			})

			Convey("Then it should wrap the underlying error", func() {
				So(err.Err, ShouldEqual, underlyingErr)
			})
		})

		Convey("When calling Error() method", func() {

			Convey("With an underlying error", func() {
				err := apperror.New("CODE", "message", 400, errors.New("cause"))

				Convey("Then it should return formatted message with cause", func() {
					So(err.Error(), ShouldEqual, "message: cause")
				})
			})

			Convey("Without an underlying error", func() {
				err := apperror.New("CODE", "message", 400, nil)

				Convey("Then it should return just the message", func() {
					So(err.Error(), ShouldEqual, "message")
				})
			})
		})

		Convey("When calling HTTPStatusCode()", func() {
			err := apperror.New("CODE", "msg", http.StatusNotFound, nil)

			Convey("Then it should return the correct HTTP status", func() {
				So(err.HTTPStatusCode(), ShouldEqual, http.StatusNotFound)
			})
		})

		Convey("When using error unwrapping", func() {
			cause := errors.New("root cause")
			err := apperror.New("CODE", "msg", 500, cause)

			Convey("Then errors.Is should find the underlying error", func() {
				So(errors.Is(err, cause), ShouldBeTrue)
			})
		})

		Convey("When using WithDetails", func() {
			err := apperror.New("CODE", "msg", 400, nil).
				WithDetails("field", "email").
				WithDetails("reason", "invalid format")

			Convey("Then details should be accessible", func() {
				So(err.Details["field"], ShouldEqual, "email")
				So(err.Details["reason"], ShouldEqual, "invalid format")
			})
		})

		Convey("When using WithError", func() {
			cause := errors.New("db connection failed")
			err := apperror.New("DB_ERROR", "database error", 500, nil).WithError(cause)

			Convey("Then the wrapped error should be accessible", func() {
				So(errors.Is(err, cause), ShouldBeTrue)
			})
		})

		Convey("When checking IsAppError", func() {
			appErr := apperror.New("CODE", "msg", 400, nil)
			regularErr := errors.New("regular error")

			Convey("Then it should identify AppError correctly", func() {
				So(apperror.IsAppError(appErr), ShouldBeTrue)
			})

			Convey("Then it should not identify regular errors", func() {
				So(apperror.IsAppError(regularErr), ShouldBeFalse)
			})
		})

		Convey("When using GetAppError", func() {
			appErr := apperror.New("CODE", "msg", 400, nil)

			Convey("Then it should extract from direct AppError", func() {
				extracted := apperror.GetAppError(appErr)
				So(extracted, ShouldNotBeNil)
				So(extracted.Code, ShouldEqual, "CODE")
			})

			Convey("Then it should return nil for non-AppError", func() {
				regularErr := errors.New("regular")
				extracted := apperror.GetAppError(regularErr)
				So(extracted, ShouldBeNil)
			})
		})
	})
}

func TestPredefinedErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		err            *apperror.AppError
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "InvalidCredentials",
			err:            apperror.InvalidCredentials(nil),
			expectedCode:   apperror.ErrCodeInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "SessionExpired",
			err:            apperror.SessionExpired(nil),
			expectedCode:   apperror.ErrCodeSessionExpired,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "NotFound",
			err:            apperror.NotFound("User", "123"),
			expectedCode:   apperror.ErrCodeNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "AlreadyExists",
			err:            apperror.AlreadyExists("Email", "test@example.com"),
			expectedCode:   apperror.ErrCodeAlreadyExists,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "ValidationFailed",
			err:            apperror.ValidationFailed("invalid input"),
			expectedCode:   apperror.ErrCodeValidationFailed,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InternalError",
			err:            apperror.InternalError(errors.New("panic")),
			expectedCode:   apperror.ErrCodeInternalError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "BusinessRuleViolation",
			err:            apperror.BusinessRuleViolation("max_habits", "Cannot create more than 10 habits"),
			expectedCode:   apperror.ErrCodeBusinessRuleViolation,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			Convey("Given "+tc.name+" error factory", t, func() {
				Convey("Then it should have the correct error code", func() {
					So(tc.err.Code, ShouldEqual, tc.expectedCode)
				})

				Convey("Then it should have the correct HTTP status", func() {
					So(tc.err.StatusCode, ShouldEqual, tc.expectedStatus)
				})
			})
		})
	}
}

func TestNotFoundDetails(t *testing.T) {
	t.Parallel()

	Convey("Given a NotFound error", t, func() {
		err := apperror.NotFound("Habit", "abc-123")

		Convey("Then it should have resource detail", func() {
			So(err.Details["resource"], ShouldEqual, "Habit")
		})

		Convey("Then it should have identifier detail", func() {
			So(err.Details["identifier"], ShouldEqual, "abc-123")
		})
	})
}
