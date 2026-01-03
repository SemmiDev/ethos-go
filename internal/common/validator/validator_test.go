package validator_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/common/validator"
)

func TestValidator(t *testing.T) {
	t.Parallel()

	Convey("Given the validator package", t, func() {

		Convey("When creating a new validator", func() {
			v := validator.New("en")

			Convey("Then it should not be nil", func() {
				So(v, ShouldNotBeNil)
			})
		})
	})
}

func TestValidateRequiredField(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string `json:"name" validate:"required"`
	}

	Convey("Given a struct with required field", t, func() {
		v := validator.New("en")

		Convey("When input is valid", func() {
			input := testStruct{Name: "John"}
			err := v.Validate(input)

			Convey("Then it should pass validation", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When required field is empty", func() {
			input := testStruct{Name: ""}
			err := v.Validate(input)

			Convey("Then it should fail validation", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	type emailStruct struct {
		Email string `json:"email" validate:"required,email"`
	}

	testCases := []struct {
		email   string
		isValid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.org", true},
		{"invalid-email", false},
		{"@nodomain.com", false},
		{"", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.email, func(t *testing.T) {
			t.Parallel()

			Convey("Given email: "+tc.email, t, func() {
				v := validator.New("en")
				input := emailStruct{Email: tc.email}
				err := v.Validate(input)

				if tc.isValid {
					Convey("Then it should pass validation", func() {
						So(err, ShouldBeNil)
					})
				} else {
					Convey("Then it should fail validation", func() {
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	}
}

func TestValidateMinMaxLength(t *testing.T) {
	t.Parallel()

	type lengthStruct struct {
		Password string `json:"password" validate:"min=6,max=20"`
	}

	testCases := []struct {
		name    string
		pwd     string
		isValid bool
	}{
		{"too short", "12345", false},
		{"minimum", "123456", true},
		{"normal", "securepassword", true},
		{"maximum", "12345678901234567890", true},
		{"too long", "123456789012345678901", false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			Convey("Given password: "+tc.name, t, func() {
				v := validator.New("en")
				input := lengthStruct{Password: tc.pwd}
				err := v.Validate(input)

				if tc.isValid {
					Convey("Then it should pass validation", func() {
						So(err, ShouldBeNil)
					})
				} else {
					Convey("Then it should fail validation", func() {
						So(err, ShouldNotBeNil)
					})
				}
			})
		})
	}
}

func TestValidationErrorsHelpers(t *testing.T) {
	t.Parallel()

	Convey("Given validation errors", t, func() {
		v := validator.New("en")

		type testStruct struct {
			Name  string `json:"name" validate:"required"`
			Email string `json:"email" validate:"required"`
		}

		Convey("When ValidateAndGetErrors is called", func() {
			errs := v.ValidateAndGetErrors(testStruct{})

			Convey("Then it should return a slice of errors", func() {
				So(len(errs), ShouldEqual, 2)
			})
		})

		Convey("When ToKV is called", func() {
			errs := v.ValidateAndGetErrors(testStruct{})
			kv := errs.ToKV()

			Convey("Then it should have name key", func() {
				_, ok := kv["name"]
				So(ok, ShouldBeTrue)
			})

			Convey("Then it should have email key", func() {
				_, ok := kv["email"]
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When Error() is called", func() {
			errs := v.ValidateAndGetErrors(testStruct{})
			errStr := errs.Error()

			Convey("Then it should return non-empty string", func() {
				So(errStr, ShouldNotBeEmpty)
			})
		})

		Convey("When IsValidationErrors is called", func() {
			err := v.Validate(testStruct{})

			Convey("Then it should return true for ValidationErrors", func() {
				So(validator.IsValidationErrors(err), ShouldBeTrue)
			})
		})

		Convey("When GetValidationErrors is called", func() {
			err := v.Validate(testStruct{})
			errs, ok := validator.GetValidationErrors(err)

			Convey("Then ok should be true", func() {
				So(ok, ShouldBeTrue)
			})

			Convey("Then errors should not be empty", func() {
				So(len(errs), ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestLocaleMessages(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string `json:"name" validate:"required"`
	}

	Convey("Given Indonesian locale", t, func() {
		v := validator.New("id")
		errs := v.ValidateAndGetErrors(testStruct{})

		Convey("Then error message should not be empty", func() {
			So(len(errs), ShouldBeGreaterThan, 0)
			So(errs[0].Message, ShouldNotBeEmpty)
		})
	})

	Convey("Given English locale", t, func() {
		v := validator.New("en")
		errs := v.ValidateAndGetErrors(testStruct{})

		Convey("Then error message should not be empty", func() {
			So(len(errs), ShouldBeGreaterThan, 0)
			So(errs[0].Message, ShouldNotBeEmpty)
		})
	})
}
