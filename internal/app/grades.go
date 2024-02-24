package app

import (
	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/KurbanowS/technical_task3/internal/school"
)

func GradeList(f models.GradeFilterRequest) ([]*models.GradeResponse, int, error) {
	grades, total, err := school.School().GradeFindBy(f)
	if err != nil {
		return nil, 0, err
	}
	gradeResponse := []*models.GradeResponse{}
	for _, grade := range grades {
		s := models.GradeResponse{}
		s.FromModel(grade)
		gradeResponse = append(gradeResponse, &s)
	}
	return gradeResponse, total, nil
}

func GradeCreate(data models.GradeRequest) (*models.GradeResponse, error) {
	model := &models.Grade{}
	data.ToModel(model)
	var err error
	model, err = school.School().GradeCreate(model)
	if err != nil {
		return nil, err
	}
	res := &models.GradeResponse{}
	res.FromModel(model)
	return res, nil
}

func GradeUpdate(data models.GradeRequest) (*models.GradeResponse, error) {
	model := &models.Grade{}
	data.ToModel(model)
	var err error
	model, err = school.School().GradeUpdate(model)
	if err != nil {
		return nil, err
	}
	res := &models.GradeResponse{}
	res.FromModel(model)
	return res, nil
}
