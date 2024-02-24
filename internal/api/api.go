package api

import (
	"net/http"

	"github.com/KurbanowS/technical_task3/internal/app"
	"github.com/KurbanowS/technical_task3/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Success(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, SuccessResponseObject(data))
}

func SuccessResponseObject(data gin.H) gin.H {
	return data
}

func BindAndValidate(c *gin.Context, r interface{}) (errMessage string, errKey string) {
	if err := c.Bind(r); err != nil {
		errMessage = err.Error()
		return
	}
	v := validator.New()
	if err := v.Struct(r); err != nil {
		err := err.(validator.ValidationErrors)[0]
		errMessage = err.Tag()
		errKey = (err.Field())
		return
	}
	return
}

func handleError(c *gin.Context, err error) {
	if errA, ok := err.(*app.AppError); ok {
		if errA == app.ErrNotFound {
			if utils.GetLoggerDesc() == "" {
				utils.Logger.Error(err)
			}
			c.JSON(http.StatusNotFound, ErrorResponseObject(errA))
		}
	}
}

func ErrorResponseObject(err *app.AppError) gin.H {
	return gin.H{
		"success": false,
		"data":    nil,
		"error": gin.H{
			"code":    err.Code(),
			"key":     err.Key(),
			"comment": err.Comment(),
		},
	}
}
