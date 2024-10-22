package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if lastErr := ctx.Errors.Last(); lastErr != nil {
			switch e := lastErr.Err.(type) {
			case *json.UnmarshalTypeError:
				ctx.AbortWithStatusJSON(apperror.ErrBadRequest.Code, dto.Response{Errors: listJSONErrors(e)})
			case validator.ValidationErrors:
				ctx.AbortWithStatusJSON(apperror.ErrBadRequest.Code, dto.Response{Errors: listValidationErrors(e)})
			case *apperror.Error:
				ctx.AbortWithStatusJSON(e.Code, dto.Response{Message: e.Error()})
			default:
				ctx.AbortWithStatusJSON(apperror.ErrInternalServerError.Code, dto.Response{Message: e.Error()})
			}
		}
	}
}

type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func listJSONErrors(typeError *json.UnmarshalTypeError) []fieldError {
	return []fieldError{
		{
			Field:   typeError.Field,
			Message: fmt.Sprintf("invalid type, expected %s", typeError.Type.String()),
		},
	}
}

func listValidationErrors(validationErrors validator.ValidationErrors) []fieldError {
	var fieldErrors []fieldError
	for _, e := range validationErrors {
		fieldErrors = append(fieldErrors, fieldError{
			Field:   strings.ToLower(e.Field()),
			Message: translate(e),
		})
	}

	return fieldErrors
}

func translate(fe validator.FieldError) string {
	switch fe.Tag() {
	case appconst.BindingRequired:
		return "this field is required"
	case appconst.BindingEq:
		return fmt.Sprintf("this field only accepts %s", fe.Param())
	case appconst.BindingEmail:
		return "this field contains invalid email format"
	case appconst.BindingMin:
		return fmt.Sprintf("this field must be at least %s characters long", fe.Param())
	case appconst.BindingGte:
		return fmt.Sprintf("this field must be greater than or equal to %s", fe.Param())
	case appconst.BindingNumeric:
		return "this field only accepts numeric characters"
	case appconst.BindingRequiredWithout:
		return fmt.Sprintf("this field is at least required if %s is empty", fe.Param())
	case appconst.BindingRequiredWithoutAll:
		return fmt.Sprintf("this field is at least required if %s is empty", fe.Param())
	case appconst.BindingOneOf:
		return fmt.Sprintf("this field only accepts %s", strings.Split(fe.Param(), ","))
	default:
		return appconst.DefaultFieldErrorTranslation
	}
}
