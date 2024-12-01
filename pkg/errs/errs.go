package errs

import (
	"net/http"
	"strconv"
)

type MessageErr interface {
	Message() string
	Status() string
	Error() string
}

type MessageErrData struct {
	ErrMessage string `json:"message"`
	ErrStatus  string `json:"status"`
	ErrError   string `json:"error"`
}

func (e *MessageErrData) Message() string {
	return e.ErrMessage
}

func (e *MessageErrData) Status() string {
	return e.ErrStatus
}

func (e *MessageErrData) Error() string {
	return e.ErrError
}

func NewCustomErrs(message string, status string, err string) MessageErr {
	return &MessageErrData{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
	}
}

func NewBadRequest(message string) MessageErr {
	return &MessageErrData{
		ErrMessage: message,
		ErrStatus:  strconv.Itoa(http.StatusBadRequest),
		ErrError:   "BAD_REQUEST",
	}
}

func NewInternalServerError(message string) MessageErr {
	return &MessageErrData{
		ErrMessage: message,
		ErrStatus:  strconv.Itoa(http.StatusInternalServerError),
		ErrError:   "INTERNAL_SERVER_ERROR",
	}
}

func NewUnprocessibleEntityError(message string) MessageErr {
	return &MessageErrData{
		ErrMessage: message,
		ErrStatus:  strconv.Itoa(http.StatusUnprocessableEntity),
		ErrError:   "INVALID_REQUEST_BODY",
	}
}

func NewUnauthorizedError(message string) MessageErr {
	return &MessageErrData{
		ErrMessage: message,
		ErrStatus:  strconv.Itoa(http.StatusUnauthorized),
		ErrError:   "UNAUTHORIZED",
	}
}
