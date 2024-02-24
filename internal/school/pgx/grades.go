package pgx

import (
	"context"
	"strconv"
	"strings"

	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const sqlGradeFields = `g.id, g.mark`
const sqlGradeInsert = `insert into grades`
const sqlGradeSelect = `select ` + sqlGradeFields + `from grades g where g.id= ANY($1::int[])`
const sqlGradeSelectMany = `select ` + sqlGradeFields + `, count(*) over() as total from grades g where g.id=g.id limit $1 offset $2`
const sqlGradeUpdate = `update grades g set id=id`

func scanGrades(rows pgx.Row, m *models.Grade, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnForScan(m, addColumns...)...)
	return
}

func (d *PgxSchool) GradeFindById(ID string) (*models.Grade, error) {
	row, err := d.GradeFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxSchool) GradeFindByIds(Ids []string) ([]*models.Grade, error) {
	grades := []*models.Grade{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlGradeSelect, (Ids))
		for rows.Next() {
			m := models.Grade{}
			err := scanGrades(rows, &m)
			if err != nil {
				return err
			}
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return grades, nil
}

func (d *PgxSchool) GradeFindBy(f models.GradeFilterRequest) (grades []*models.Grade, total int, err error) {
	args := []interface{}{f.Limit, f.Offset}
	qs, args := GradeListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), qs, args...)
		for rows.Next() {
			grade := models.Grade{}
			err = scanGrades(rows, &grade, &total)
			if err != nil {
				return err
			}
			grades = append(grades, &grade)
		}
		return
	})
	if err != nil {
		return nil, 0, err
	}
	return grades, total, nil
}

func (d *PgxSchool) GradeCreate(model *models.Grade) (*models.Grade, error) {
	qs, args := GradeCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		err = tx.QueryRow(context.Background(), qs, args...).Scan(&model.ID)
		return
	})
	if err != nil {
		return nil, err
	}
	editModel, err := d.GradeFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxSchool) GradeUpdate(model *models.Grade) (*models.Grade, error) {
	qs, args := GradeUpdateQuery(model)
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		return nil, err
	}
	editModel, err := d.GradeFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func GradeCreateQuery(m *models.Grade) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := GradeAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		k += ", " + cols
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlGradeInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func GradeUpdateQuery(m *models.Grade) (string, []interface{}) {
	args := []interface{}{}
	sets := ""
	q := GradeAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		sets = ", " + k + "=$" + strconv.Itoa(len(args))
	}
	args = append(args, m.ID)
	qs := strings.ReplaceAll(sqlGradeUpdate, "set id=id", "set id=id "+sets+" ") + " where id=$" + strconv.Itoa(len(args))
	return qs, args
}

func GradeAtomicQuery(m *models.Grade) map[string]interface{} {
	q := map[string]interface{}{}
	q["mark"] = m.Mark
	return q
}

func GradeListBuildQuery(f models.GradeFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""
	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += "and g.id=$" + strconv.Itoa(len(args))
	}
	wheres += "order by g.id desc"
	qs := sqlGradeSelectMany
	qs = strings.ReplaceAll(qs, "g.id=g.id", "g.id=g.id"+wheres+" ")
	return qs, args
}
