package postgres

import (
	"banner/internal/db"
	"context"
	"fmt"
	"strings"
	"testing"
)

type TDB struct {
	DB *db.Database
}

func NewDB(dsn string) *TDB {
	db, err := db.NewDB(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	//d.DB.BeginTX()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) TearDown(tableName ...string) {
	//d.DB.RollBack()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {

	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY", strings.Join(tableName, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
