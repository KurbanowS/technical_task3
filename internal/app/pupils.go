package app

import (
	"errors"
	"strings"

	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/KurbanowS/technical_task3/internal/school"
)

func PupilList(f models.PupilFilterRequest) ([]*models.PupilResponse, int, error) {
	pupils, total, err := school.School().PupilFindBy(f)
	if err != nil {
		return nil, 0, err
	}
	err = school.School().PupilLoadRelations(&pupils)
	if err != nil {
		return nil, 0, err
	}

	pupilResponse := []*models.PupilResponse{}
	for _, pupil := range pupils {
		s := models.PupilResponse{}
		s.FromModel(pupil)
		pupilResponse = append(pupilResponse, &s)
	}
	return pupilResponse, total, err
}

func PupilCreate(data models.PupilRequest) (*models.PupilResponse, error) {
	model := &models.Pupil{}
	data.ToModel(model)
	res := &models.PupilResponse{}
	var err error
	model, err = school.School().PupilCreate(model)
	if err != nil {
		return nil, err
	}
	err = school.School().PupilLoadRelations(&[]*models.Pupil{model})
	if err != nil {
		return nil, err
	}

	res.FromModel(model)
	return res, nil
}

func PupilDelete(ids []string) ([]*models.PupilResponse, error) {
	pupils, err := school.School().PupilFindByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(pupils) < 1 {
		return nil, errors.New("model not found: " + strings.Join(ids, ","))
	}
	pupils, err = school.School().PupilDelete(pupils)
	if err != nil {
		return nil, err
	}
	if len(pupils) == 0 {
		return make([]*models.PupilResponse, 0), nil
	}
	var pupilResponse = []*models.PupilResponse{}
	for _, pupil := range pupils {
		var p models.PupilResponse
		p.FromModel(pupil)
		pupilResponse = append(pupilResponse, &p)
	}
	return pupilResponse, nil
}
