package user_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/random"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	Convey("Given the User domain", t, func() {

		Convey("When creating a new user with email auth", func() {
			userID := random.NewUUID()
			u := user.NewUser(userID, "test@example.com", "Test User", "hashedpassword123")

			Convey("Then it should have correct UserID", func() {
				So(u.UserID(), ShouldEqual, userID)
			})

			Convey("Then it should have correct Email", func() {
				So(u.Email(), ShouldEqual, "test@example.com")
			})

			Convey("Then it should have correct Name", func() {
				So(u.Name(), ShouldEqual, "Test User")
			})

			Convey("Then it should have hashed password", func() {
				So(*u.HashedPassword(), ShouldEqual, "hashedpassword123")
			})

			Convey("Then auth provider should be 'email'", func() {
				So(u.AuthProvider(), ShouldEqual, "email")
			})

			Convey("Then it should be active", func() {
				So(u.IsActive(), ShouldBeTrue)
			})

			Convey("Then it should NOT be verified", func() {
				So(u.IsVerified(), ShouldBeFalse)
			})

			Convey("Then it should have default timezone", func() {
				So(u.Timezone(), ShouldEqual, "Asia/Jakarta")
			})
		})
	})
}

func TestNewGoogleUser(t *testing.T) {
	t.Parallel()

	Convey("Given Google OAuth registration", t, func() {

		Convey("When creating a new Google user", func() {
			userID := random.NewUUID()
			u := user.NewGoogleUser(userID, "google@example.com", "Google User", "google-id-12345")

			Convey("Then it should have correct UserID", func() {
				So(u.UserID(), ShouldEqual, userID)
			})

			Convey("Then it should have correct Email", func() {
				So(u.Email(), ShouldEqual, "google@example.com")
			})

			Convey("Then it should have correct Name", func() {
				So(u.Name(), ShouldEqual, "Google User")
			})

			Convey("Then hashed password should be nil", func() {
				So(u.HashedPassword(), ShouldBeNil)
			})

			Convey("Then auth provider should be 'google'", func() {
				So(u.AuthProvider(), ShouldEqual, "google")
			})

			Convey("Then auth provider ID should be set", func() {
				So(*u.AuthProviderID(), ShouldEqual, "google-id-12345")
			})

			Convey("Then it should be verified (implicit for Google)", func() {
				So(u.IsVerified(), ShouldBeTrue)
			})

			Convey("Then it should be active", func() {
				So(u.IsActive(), ShouldBeTrue)
			})
		})
	})
}
