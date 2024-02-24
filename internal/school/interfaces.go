package school

import "github.com/KurbanowS/technical_task3/internal/models"

type ISchool interface {
	PupilFindById(ID string) (*models.Pupil, error)
	PupilFindByIds(IDs []string) ([]*models.Pupil, error)
	PupilFindBy(f models.PupilFilterRequest) (pupils []*models.Pupil, total int, err error)
	PupilCreate(model *models.Pupil) (*models.Pupil, error)
	PupilDelete(items []*models.Pupil) ([]*models.Pupil, error)
	PupilLoadRelations(l *[]*models.Pupil) error

	GradeFindById(ID string) (*models.Grade, error)
	GradeFindByIds(Ids []string) ([]*models.Grade, error)
	GradeFindBy(f models.GradeFilterRequest) (grades []*models.Grade, total int, err error)
	GradeCreate(model *models.Grade) (*models.Grade, error)
	GradeUpdate(model *models.Grade) (*models.Grade, error)
}
