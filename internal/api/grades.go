package api

import (
	"strconv"

	"github.com/KurbanowS/technical_task3/internal/app"
	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/gin-gonic/gin"
)

func GradeRoutes(api *gin.RouterGroup) {
	gradeRoutes := api.Group("/grades")
	{
		gradeRoutes.GET("", GradeList)
		gradeRoutes.POST("", GradeCreate)
		gradeRoutes.PUT(":id", GradeUpdate)
	}
}

func GradeList(c *gin.Context) {
	r := models.GradeFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	grades, total, err := app.GradeList(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"grades": grades,
		"total":  total,
	})
}

func GradeCreate(c *gin.Context) {
	r := models.GradeRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	grades, err := app.GradeCreate(r)
	if err != nil {
		handleError(c, err)
	}
	Success(c, gin.H{
		"grades": grades,
	})
}

func GradeUpdate(c *gin.Context) {
	r := models.GradeRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	idp := uint(id)
	r.ID = &idp

	if id == 0 {
		handleError(c, app.ErrRequired.SetKey("id"))
		return
	}
	grade, err := app.GradeUpdate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"grade_update": grade,
	})
}
