package main

import (
	"log"

	"github.com/KurbanowS/technical_task3/config"
	"github.com/KurbanowS/technical_task3/internal/api"
	"github.com/KurbanowS/technical_task3/internal/school"
	"github.com/KurbanowS/technical_task3/internal/school/pgx"
	"github.com/KurbanowS/technical_task3/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	defer utils.InitLogs().Close()
	config.LoadConfig()
	defer school.Init().(*pgx.PgxSchool).Close()

	routes := gin.Default()
	api.Routes(routes)
	if err := routes.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
