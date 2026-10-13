package httpResponse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"spotsync-backend/internal/apperror"
)

// Validator plugs go-playground/validator into Echo (c.Validate(&req)).
type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	// Report JSON field names ("license_plate"), not Go names ("LicensePlate").
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})
	return &Validator{v: v}
}

func (cv *Validator) Validate(i any) error {
	err := cv.v.Struct(i)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		return apperror.Invalid("Validation failed").WithDetails(formatErrors(verrs))
	}
	return err
}

func formatErrors(verrs validator.ValidationErrors) string {
	parts := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		parts = append(parts, fe.Field()+" "+describe(fe))
	}
	return strings.Join(parts, "; ")
}

func describe(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters", fe.Param())
		}
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "is invalid"
	}
}
