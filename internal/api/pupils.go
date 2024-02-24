package api

import (
	"github.com/KurbanowS/technical_task3/internal/app"
	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/gin-gonic/gin"
)

func PupilRoutes(api *gin.RouterGroup) {
	pupilRoutes := api.Group("/pupils")
	{
		pupilRoutes.GET("", PupilList)
		pupilRoutes.POST("", PupilCreate)
		pupilRoutes.DELETE("", PupilDelete)
	}
}

func PupilList(c *gin.Context) {
	r := models.PupilFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	pupils, total, err := app.PupilList(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"pupils": pupils,
		"total":  total,
	})
}

func PupilCreate(c *gin.Context) {
	r := models.PupilRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	pupil, err := app.PupilCreate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"pupil": pupil,
	})
}

func PupilDelete(c *gin.Context) {
	var ids []string = c.QueryArray("ids")
	if len(ids) == 0 {
		handleError(c, app.ErrRequired.SetKey("ids"))
		return
	}
	pupils, err := app.PupilDelete(ids)
	if err != nil {
		handleError(c, err)
	}
	Success(c, gin.H{
		"pupils": pupils,
	})
}
