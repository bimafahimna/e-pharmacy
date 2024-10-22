package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrUniqueViolation = errors.New("unique violation")
)

var (
	ErrNoChangesMade                = newError(http.StatusOK, "No changes made")
	ErrIDMustBeNumeric              = newError(http.StatusBadRequest, "ID must be numeric")
	ErrInvalidPasswordFormat        = newError(http.StatusBadRequest, "Password should contain at least a lowercase, an uppercase, a digit, and a special character")
	ErrEmailHasBeenRegistered       = newError(http.StatusBadRequest, "Email has been registered")
	ErrPartnerNameHasBeenRegistered = newError(http.StatusBadRequest, "Partner name has been registered")
	ErrProductNameHasBeenRegistered = newError(http.StatusBadRequest, "Cannot add the same product with the same manufacturer, name, and generic name")
	ErrInvalidProductClassification = newError(http.StatusBadRequest, "Invalid product classification id")
	ErrInvalidProductForm           = newError(http.StatusBadRequest, "Invalid product form id")
	ErrInvalidManufacturer          = newError(http.StatusBadRequest, "Invalid manufacturer id")
	ErrVerifyRegisteredEmail        = newError(http.StatusBadRequest, "You are registered, please check your email to verify your account")
	ErrInvalidVerificationToken     = newError(http.StatusBadRequest, "Invalid verification token")
	ErrPasswordDoNotMatch           = newError(http.StatusBadRequest, "Password do not match")
	ErrEmailNotFound                = newError(http.StatusBadRequest, "Email is not registered")
	ErrIncorrectEmailOrPassword     = newError(http.StatusBadRequest, "Incorrect email or password")
	ErrInvalidPasswordResetToken    = newError(http.StatusBadRequest, "Invalid password reset token")
	ErrResetPasswordTokenExists     = newError(http.StatusBadRequest, "Please check your email to reset your password")
	ErrCannotResetPassword          = newError(http.StatusBadRequest, "Cannot reset password for Google logged in users")
	ErrCodeTokenExchangeFailed      = newError(http.StatusBadRequest, "Code-Token exchange failed")
	ErrBadRequest                   = newError(http.StatusBadRequest, "Bad request")
	ErrInvalidSort                  = newError(http.StatusBadRequest, "Invalid sort")
	ErrInvalidDayId                 = newError(http.StatusBadRequest, "Invalid day id")
	ErrEmptyActiveDays              = newError(http.StatusBadRequest, "Active days cannot be empty")
	ErrInvalidHourFormat            = newError(http.StatusBadRequest, "Invalid hour format")
	ErrAddressNameAlreadyExist      = newError(http.StatusBadRequest, "Address name is already exist")
	ErrPharmacyLocationAlreadyExist = newError(http.StatusBadRequest, "Pharmacy at this location is already exist")
	ErrPharmacyProductUnavailable   = newError(http.StatusBadRequest, "Unable to add unavailable product")
	ErrPharmacyProductHasBeenBought = newError(http.StatusBadRequest, "Unable to delete product that has been bought by user")
	ErrPaymentNotFound              = newError(http.StatusBadRequest, "Payment not found")
	ErrFileNotFound                 = newError(http.StatusBadRequest, "File not found")
	ErrFileSizeExceededLimit        = newError(http.StatusBadRequest, "File size exceeded 500kb limit")
	ErrUnsupportedFileType          = newError(http.StatusBadRequest, "Unsupported file type")
	ErrUnavailableLogistic          = newError(http.StatusBadRequest, "Unavailable logistic")
	ErrPharmacistMustBeUnique       = newError(http.StatusConflict, "SIPA and WhatsApp number must be unique")
	ErrProductCategoryMustBeUnique  = newError(http.StatusConflict, "Product category must be unique")
	ErrPharmacyProductAlreadyExists = newError(http.StatusConflict, "Pharmacy product already exists")
	ErrUpdatedPharmacyProductToday  = newError(http.StatusBadRequest, "Pharmacy product can only be updated once")
	ErrPharmacistIsAlreadyAssigned  = newError(http.StatusBadRequest, "Pharmacist is alreadys assigned")
	ErrInvalidToken                 = newError(http.StatusUnauthorized, "Invalid token")
	ErrUnauthorized                 = newError(http.StatusUnauthorized, "Unauthorized")
	ErrForbiddenRole                = newError(http.StatusForbidden, "Forbidden role")
	ErrForbiddenNotVerified         = newError(http.StatusForbidden, "Not verified")
	ErrUserFetchFailed              = newError(http.StatusServiceUnavailable, "User data fetch failed")
	ErrNotFound                     = newError(http.StatusNotFound, "Resource not found")
	ErrInternalServerError          = newError(http.StatusInternalServerError, "Internal server error")
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func newError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
