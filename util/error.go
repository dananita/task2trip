package util

import (
	"github.com/go-pg/pg"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/go-errors/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code        int    `json:"-"`
	ErrorKey    string `json:"key"`
	Description string `json:"description"`
	Inner       error  `json:"inner"`
	Stacktrace  string `json:"stacktrace"`
}

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return he.Description
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, key, description string, inner ...error) *HTTPError {
	he := &HTTPError{
		Code:        code,
		ErrorKey:    key,
		Description: description,
	}
	if len(inner) > 0 {
		he.Inner = inner[0]
		he.Description = he.Description + ": " + he.Inner.Error()
		if withTrace, ok := inner[0].(*errors.Error); ok {
			he.Stacktrace = withTrace.ErrorStack()
		}
	}
	if len(he.Stacktrace) == 0 {
		he.Stacktrace = string(debug.Stack())
	}

	return he
}

func (he *HTTPError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	if he.Code > 0 {
		rw.WriteHeader(he.Code)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	if err := producer.Produce(rw, he); err != nil {
		panic(err)
	}
}

func NewNotFoundError(key, description string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, key, description)
}

func NewUnauthorized(key, description string) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, key, description)
}

func NewBadRequestError(key, description string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, key, description)
}

func NewServerError(key, description string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, key, description)
}

func ConvertHTTPErrorToResponse(err error) middleware.Responder {
	if httpError, ok := err.(*HTTPError); ok {
		return httpError
	}
	if err == pg.ErrNoRows {
		return NewNotFoundError("not found", "")
	}

	return NewServerError("internal.server.error", err.Error())
}

func ErrWithTrace(err error, optionalAnnotation ...string) error {
	annotation := ""
	if len(optionalAnnotation) > 0 {
		annotation = optionalAnnotation[0]
	}
	if err == nil {
		return nil
	}

	if len(annotation) > 0 {
		return errors.WrapPrefix(err, annotation, 1)
	}
	return errors.New(err)
}

// NewBadRequestExtendedError indicates that client send invalid request
func NewBadRequestExtendedError(key, description string) error {
	return NewHTTPError(http.StatusBadRequest, key, description, nil)
}

func MultiError(errs ...error) error {
	var notNil []error
	for _, err := range errs {
		if err != nil {
			notNil = append(notNil, err)
		}
	}
	if len(notNil) == 0 {
		return nil
	}
	if len(notNil) == 1 {
		return notNil[0]
	}
	// var errStrings []string
	errStrings := make([]string, 0)
	for _, err := range notNil {
		errStrings = append(errStrings, err.Error())
	}

	return errors.Errorf("Multiple errors: %s", strings.Join(errStrings, ", "))
}
