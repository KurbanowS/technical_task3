package pgx

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/KurbanowS/technical_task3/config"
	"github.com/KurbanowS/technical_task3/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxSchool struct {
	pool *pgxpool.Pool
}

func (d PgxSchool) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *PgxSchool) Close() {
	d.pool.Close()
}

func Init() *PgxSchool {
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable connect_timeout=5", config.Conf.DbUsername, config.Conf.DbDatabase, config.Conf.DbPassword, config.Conf.DbHost, config.Conf.DbPort)
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		log.Fatal(err)
	}
	return &PgxSchool{pool: pool}
}

type pgxQuery func(conn *pgxpool.Conn) (err error)

func (d *PgxSchool) runQuery(ctx context.Context, f pgxQuery) (err error) {
	err = d.Pool().AcquireFunc(ctx, f)
	if err != nil {
		return err
	}
	return
}

func parseColumnForScan(sub interface{}, addColumns ...interface{}) []interface{} {
	s := reflect.ValueOf(sub).Elem()
	numCols := s.NumField() - len(sub.(models.HasRelationFields).RelationFields())
	columns := []interface{}{}
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns = append(columns, field.Addr().Interface())
	}
	columns = append(columns, addColumns...)
	return columns
}
