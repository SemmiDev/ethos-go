package random

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/semmidev/ethos-go/internal/common/apperror"
)

const (
	minOTPLength = 4
	maxOTPLength = 12
	digits       = "0123456789"
)

// GenerateNumericOTP menghasilkan OTP numerik (digit-only) dengan panjang tertentu.
// Panjang valid antara minOTPLength dan maxOTPLength.
func GenerateNumericOTP(length int) (string, error) {
	if length < minOTPLength || length > maxOTPLength {
		return "", apperror.ValidationFailed("OTP length must be between 4 and 12 digits")
	}

	digitsLength := big.NewInt(int64(len(digits)))
	var otpBuilder strings.Builder
	otpBuilder.Grow(length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, digitsLength)
		if err != nil {
			return "", apperror.InternalError(err)
		}
		otpBuilder.WriteByte(digits[n.Int64()])
	}

	return otpBuilder.String(), nil
}
