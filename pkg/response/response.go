package response

import (
	"net/http"

	"repo-backend/pkg/apperr"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success   bool       `json:"success"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data, RequestID: requestID(c)})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data, RequestID: requestID(c)})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, err error) {
	appError := apperr.As(err)
	c.JSON(appError.HTTPStatus, Envelope{
		Success:   false,
		Error:     &ErrorBody{Code: appError.Code, Message: appError.Message},
		RequestID: requestID(c),
	})
}

func requestID(c *gin.Context) string {
	return c.GetString("request_id")
}
