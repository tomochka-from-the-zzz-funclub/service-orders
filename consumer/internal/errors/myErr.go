package myErrors

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

type Error struct {
	httpCode int
	cause    string
}

func NewError(code int, cause string) Error {
	return Error{
		httpCode: code,
		cause:    cause,
	}
}

func (e Error) GetHttpCode() int {
	return e.httpCode
}

func (e Error) GetCause() string {
	return e.cause
}
func (e Error) Error() string {
	return "Status code: " + strconv.Itoa(e.httpCode) + " cause: " + e.cause
}

// json
var ErrParseJSON = NewError(fasthttp.StatusBadRequest, "error decoding json")
var ErrEqualJSON = NewError(fasthttp.StatusBadRequest, "error read information in JSON format: empty")

// database
var ErrCreatePostgresConnection = NewError(fasthttp.StatusInternalServerError, "don't create postgres connection")
var ErrPing = NewError(fasthttp.StatusInternalServerError, "error ping postgres")
var ErrNotFoundOrder = NewError(fasthttp.StatusNotFound, "not found order in bd")

// http
var ErrMethodNotAllowed = NewError(fasthttp.StatusMethodNotAllowed, "method not allowed")

// for test service
var ErrNotFoundOrderCache = NewError(fasthttp.StatusNotFound, "not found order in cache")
