package api

import (
	"github.com/KurbanowS/technical_task3/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes(routes *gin.Engine) {
	if config.Conf.AppEnvIsProd {
		gin.SetMode("release")
	}

	routes.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
	}))

	api := routes.Group("/api")
	{
		GradeRoutes(api)
		PupilRoutes(api)
	}
}
