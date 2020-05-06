package fcerr

import (
	"fmt"
	"net/http"
)

//FCErr is a custom Error interface that uses a message and an int to provide consistant, informative error support.
type FCErr interface {
	Message() string
	Status() int
	Error() string
}

type fcerr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

//Message is the simple message to return the string that comprises the message of the error.
func (e fcerr) Message() string {
	return e.ErrMessage
}

func (e fcerr) Status() int {
	return e.ErrStatus
}

func (e fcerr) Error() string {
	return fmt.Sprint("Message: ", e.ErrMessage, " - Status: ", e.ErrStatus)
}

//NewFCErr takes a message string and a status int and gives you the final FCErr object.
func NewFCErr(message string, status int) FCErr {
	err := fmt.Sprint("Message: ", message, " - Status: ", status)
	return fcerr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
	}
}

//NewInternalServerError takes a message string and gives you a FCErr object with the status of http.StatusInternalServerError.
func NewInternalServerError(message string) FCErr {
	err := fmt.Sprint("Message: ", message, " - Status: ", http.StatusInternalServerError)
	return fcerr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   err,
	}
}

//NewBadRequestError takes a message string and gives you a FCErr object with the status of http.StatusInternalServerError.
func NewBadRequestError(message string) FCErr {
	err := fmt.Sprint("Message: ", message, " - Status: ", http.StatusBadRequest)
	return fcerr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   err,
	}
}

//NewNotFoundError takes a message string and gives you a FCErr object with the status of http.StatusNotFound.
func NewNotFoundError(message string) FCErr {
	err := fmt.Sprint("Message: ", message, " - Status: ", http.StatusNotFound)
	return fcerr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   err,
	}
}
