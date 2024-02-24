package pgx

import (
	"context"
	"strconv"
	"strings"

	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const sqlPupilFields = `p.id, p.name, p.grade_id`
const sqlPupilInsert = `insert into pupils`
const sqlPupilSelect = `select ` + sqlPupilFields + ` from pupils p where p.id = ANY($1::int[])`
const sqlPupilSelectMany = `select ` + sqlPupilFields + `, count(*) over() as total from pupils p where p.id=p.id limit $1 offset $2`
const sqlPupilDelete = `delete from pupils p where id = ANY[$1::int[]]`
const sqlPupilGrade = `select ` + sqlGradeFields + `, p.id from products p right join grades g on (g.id=p.grade_id) where p.id = ANY($1::int[])`

func scanPupils(rows pgx.Row, m *models.Pupil, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnForScan(m, addColumns...)...)
	return
}

func (d *PgxSchool) PupilFindById(ID string) (*models.Pupil, error) {
	row, err := d.PupilFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxSchool) PupilFindByIds(Ids []string) ([]*models.Pupil, error) {
	pupils := []*models.Pupil{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlPupilSelect, (Ids))
		for rows.Next() {
			m := models.Pupil{}
			err := scanPupils(rows, &m)
			if err != nil {
				return err
			}
			pupils = append(pupils, &m)
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return pupils, nil
}

func (d *PgxSchool) PupilFindBy(f models.PupilFilterRequest) (pupils []*models.Pupil, total int, err error) {
	args := []interface{}{f.Limit, f.Offset}
	qs, args := PupilListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), qs, args...)
		for rows.Next() {
			pupil := models.Pupil{}
			err = scanPupils(rows, &pupil, &total)
			if err != nil {
				return err
			}
			pupils = append(pupils, &pupil)
		}
		return
	})
	if err != nil {
		return nil, 0, err
	}
	return pupils, total, nil
}

func (d *PgxSchool) PupilCreate(model *models.Pupil) (*models.Pupil, error) {
	qs, args := PupilCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		err = tx.QueryRow(context.Background(), qs, args...).Scan()
		return
	})
	if err != nil {
		return nil, err
	}
	editModel, err := d.PupilFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxSchool) PupilDelete(items []*models.Pupil) ([]*models.Pupil, error) {
	ids := []uint{}
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), sqlPupilDelete, (ids))
		return
	})
	if err != nil {
		return nil, err
	}
	return items, nil

}

func PupilCreateQuery(m *models.Pupil) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := PupilAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		cols += ", " + k
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlPupilInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func PupilAtomicQuery(m *models.Pupil) map[string]interface{} {
	q := map[string]interface{}{}
	q["name"] = m.Name
	q["grade_id"] = m.GradeId
	return q
}

func PupilListBuildQuery(f models.PupilFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""
	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += " and p.id=$" + strconv.Itoa(len(args))
	}
	if f.GradeId != nil && *f.GradeId != 0 {
		args = append(args, *f.GradeId)
		wheres += " and p.grade_id=$" + strconv.Itoa(len(args))
	}
	wheres += "order by p.id desc"
	qs := sqlPupilSelectMany
	qs = strings.ReplaceAll(qs, "p.id=p.id", "p.id=p.id "+wheres)
	return qs, args
}

func (d *PgxSchool) PupilLoadRelations(l *[]*models.Pupil) error {
	ids := []string{}
	for _, m := range *l {
		ids = append(ids, strconv.Itoa(int(m.ID)))
	}
	if len(ids) < 1 {
		return nil
	}

	if rs, err := d.PupilLoadGrade(ids); err != nil {
		return err
	} else {
		for _, r := range rs {
			for _, m := range *l {
				if r.Id == m.ID {
					m.Grade = r.Relation
				}
			}
		}
	}
	return nil

}

type PupilLoadGradeItem struct {
	Id       uint
	Relation *models.Grade
}

func (d *PgxSchool) PupilLoadGrade(ids []string) ([]PupilLoadGradeItem, error) {
	res := []PupilLoadGradeItem{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlPupilGrade, (ids))
		for rows.Next() {
			sub := models.Grade{}
			pid := uint(0)
			err = scanGrades(rows, &sub, &pid)
			if err != nil {
				return err
			}
			res = append(res, PupilLoadGradeItem{Id: pid, Relation: &sub})
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
