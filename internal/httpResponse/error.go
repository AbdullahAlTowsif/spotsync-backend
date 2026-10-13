package httpResponse

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"spotsync-backend/internal/apperror"
)

const (
	internalMessage = "Internal server error"
	internalDetails = "unexpected server error"
)

type ErrorBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  string `json:"errors,omitempty"`
}

// ErrorHandler is Echo's centralized error handler. Every error returned by a
// handler or middleware ends up here and is converted to the standard format.
// Note: in Echo v5 the signature is (c, err); v4 had them the other way round.
func ErrorHandler(c *echo.Context, err error) {
	// If the response has already started we can't change it anymore.
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil && resp.Committed {
		return
	}

	status, message, details := resolve(err)

	if status >= http.StatusInternalServerError {
		c.Logger().Error("request failed",
			"error", err,
			"method", c.Request().Method,
			"path", c.Request().URL.Path,
		)
	}

	_ = c.JSON(status, ErrorBody{Success: false, Message: message, Errors: details})
}

// resolve turns any error into (status, message, details) that is safe to
// show to the client. Unknown errors become a generic 500: raw GORM/driver
// errors never leave the server.
func resolve(err error) (int, string, string) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		status := statusFor(appErr.Kind)
		if status >= http.StatusInternalServerError {
			return status, internalMessage, internalDetails
		}
		return status, appErr.Message, appErr.Details
	}

	// Errors produced by Echo itself: unknown route (404), wrong method (405),
	// malformed JSON body (400), etc.
	if code := echo.StatusCode(err); code != 0 {
		if code >= http.StatusInternalServerError {
			return code, internalMessage, internalDetails
		}
		details := ""
		var he *echo.HTTPError
		if errors.As(err, &he) && he.Message != http.StatusText(code) {
			details = he.Message
		}
		return code, http.StatusText(code), details
	}

	return http.StatusInternalServerError, internalMessage, internalDetails
}

func statusFor(kind apperror.Kind) int {
	switch kind {
	case apperror.KindInvalid:
		return http.StatusBadRequest
	case apperror.KindUnauthorized:
		return http.StatusUnauthorized
	case apperror.KindForbidden:
		return http.StatusForbidden
	case apperror.KindNotFound:
		return http.StatusNotFound
	case apperror.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
