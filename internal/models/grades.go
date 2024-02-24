package models

type Grade struct {
	ID   uint   `json:"id"`
	Mark *uint8 `json:"mark"`
}

func (Grade) HasRelationFields() []string {
	return []string{}
}

type GradeRequest struct {
	ID   *uint  `json:"id" form:"id"`
	Mark *uint8 `json:"mark" form:"mark"`
}

type GradeResponse struct {
	ID   uint   `json:"id"`
	Mark *uint8 `json:"mark"`
}

func (b *GradeRequest) ToModel(m *Grade) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Mark = b.Mark
}

func (r *GradeResponse) FromModel(m *Grade) {
	r.ID = m.ID
	r.Mark = m.Mark
}

type GradeFilterRequest struct {
	ID *uint `form:"id"`
	PaginationRequest
}
