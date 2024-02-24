package school

import "github.com/KurbanowS/technical_task3/internal/school/pgx"

var school ISchool

func School() ISchool {
	return school
}

func Init() ISchool {
	school = pgx.Init()
	return school
}
