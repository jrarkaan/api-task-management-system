package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	statusSuccess uint16 = 1
	statusError   uint16 = 0
)

func StatusOK(ctx *gin.Context, meta interface{}, message string) {
	ctx.JSON(http.StatusOK, Response{
		Meta:    meta,
		Message: message,
		Status:  statusSuccess,
		Data:    gin.H{},
	})
}

func StatusCreated(ctx *gin.Context, meta interface{}, message string) {
	ctx.JSON(http.StatusCreated, Response{
		Meta:    meta,
		Message: message,
		Status:  statusSuccess,
		Data:    gin.H{},
	})
}

func Success(ctx *gin.Context, meta interface{}, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Meta:    meta,
		Message: "Success",
		Status:  statusSuccess,
		Data:    data,
	})
}

func Created(ctx *gin.Context, meta interface{}, data interface{}, message string) {
	if message == "" {
		message = "Created successfully"
	}

	ctx.JSON(http.StatusCreated, Response{
		Meta:    meta,
		Message: message,
		Status:  statusSuccess,
		Data:    data,
	})
}

func BadRequest(ctx *gin.Context, errorMessage string) {
	errorResponse(ctx, http.StatusBadRequest, errorMessage)
}

func UnAuthorized(ctx *gin.Context, meta interface{}, errorMessage string) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Meta:    meta,
		Message: "",
		Status:  statusError,
		Error:   newXError(http.StatusUnauthorized, errorMessage),
	})
}

func Forbidden(ctx *gin.Context, meta interface{}, errorMessage string) {
	ctx.JSON(http.StatusForbidden, Response{
		Meta:    meta,
		Message: "",
		Status:  statusError,
		Error:   newXError(http.StatusForbidden, errorMessage),
	})
}

func DataNotFound(ctx *gin.Context) {
	errorResponse(ctx, http.StatusNotFound, "data not found")
}

func NotFound(ctx *gin.Context, meta interface{}, errorMessage string) {
	ctx.JSON(http.StatusNotFound, Response{
		Meta:    meta,
		Message: "",
		Status:  statusError,
		Error:   newXError(http.StatusNotFound, errorMessage),
	})
}

func Conflict(ctx *gin.Context, meta interface{}, errorMessage string) {
	ctx.JSON(http.StatusConflict, Response{
		Meta:    meta,
		Message: "",
		Status:  statusError,
		Error:   newXError(http.StatusConflict, errorMessage),
	})
}

func ServerError(ctx *gin.Context, errorMessage string) {
	errorResponse(ctx, http.StatusInternalServerError, errorMessage)
}

func ResponseError(ctx *gin.Context, meta interface{}, statusCode int, err error) {
	message := "internal server error"
	if err != nil {
		message = err.Error()
	}

	ctx.JSON(statusCode, Response{
		Meta:    meta,
		Message: "",
		Status:  statusError,
		Error:   newXError(statusCode, message),
	})
}

func errorResponse(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, Response{
		Message: "",
		Status:  statusError,
		Error:   newXError(statusCode, errorMessage),
	})
}

func newXError(code int, message string) XError {
	return XError{
		Code:    uint16(code),
		Message: message,
		Status:  true,
	}
}
