package models

type Pupil struct {
	ID      uint    `json:"id"`
	Name    *string `json:"name"`
	GradeId *uint   `json:"grade_id"`
	Grade   *Grade  `json:"grade"`
}

func (Pupil) RelationFields() []string {
	return []string{"Grade"}
}

type PupilRequest struct {
	ID      *uint   `json:"id" form:"id"`
	Name    *string `json:"name" form:"name"`
	GradeId *uint   `json:"grade_id" form:"grade_id"`
}

type PupilResponse struct {
	ID    uint           `json:"id"`
	Name  *string        `json:"name"`
	Grade *GradeResponse `json:"grade"`
}

func (b *PupilRequest) ToModel(m *Pupil) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Name = b.Name
	m.GradeId = b.GradeId
}

func (r *PupilResponse) FromModel(m *Pupil) {
	r.ID = m.ID
	r.Name = m.Name
	if r.Grade != nil {
		r.Grade = &GradeResponse{}
		r.Grade.FromModel(m.Grade)
	}
}

type PupilFilterRequest struct {
	ID      *uint `form:"id"`
	GradeId *uint `form:"grade_id"`
	PaginationRequest
}
