package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pakohan/go/sqlrepo"
)

type modelInfo struct {
	Model   Model
	Schema  string
	Columns []column
}

type column struct {
	Name       string `db:"name"`
	DataType   string `db:"data_type"`
	IsNullable bool   `db:"is_nullable"`
}

func getModelInfo(db *sqlx.DB, q *sqlrepo.SQLRepository, schema string, model Model) (*modelInfo, error) {
	columns := []column{}
	err := db.Select(&columns, q.Query("get_columns"), schema, model.TableName)
	if err != nil {
		return nil, err
	}

	return &modelInfo{
		Model:   model,
		Schema:  schema,
		Columns: columns,
	}, nil
}

func getTables(db *sqlx.DB, q *sqlrepo.SQLRepository, schema string) ([]Model, error) {
	tables := []Model{}
	err := db.Select(&tables, q.Query("get_tables"), schema)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
