package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var malformedJSON = "MalformedJSON"

type HTTPError struct {
	statusCode int

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(err error) HTTPError {

	var serr *json.SyntaxError
	if errors.As(err, &serr) {
		return HTTPError{
			statusCode: http.StatusBadRequest,
			Code:       "malformedjson",
			Message:    "Bad Request",
		}
	}

	switch err {
	case io.EOF:
		return HTTPError{
			statusCode: http.StatusBadRequest,
			Code:       "eof",
			Message:    "EOF reading HTTP request body",
		}
	case sql.ErrNoRows:
		return HTTPError{
			statusCode: http.StatusNotFound,
			Code:       "not_found",
			Message:    "Not Found",
		}
	}

	return HTTPError{
		statusCode: http.StatusInternalServerError,

		Code:    "internal",
		Message: "Internal server error",
	}
}

func (srv *Server) BadRequest(err error) HTTPError {
	return HTTPError{
		statusCode: http.StatusBadRequest,
		Code:       "bad_request",
		Message:    err.Error(),
	}
}

func (srv *Server) UnexpectedError(err error) HTTPError {
	msg := fmt.Sprintf("Unexpected error: %s", err.Error())
	srv.logger.PrintError(msg)

	return NewHTTPError(err)
}
