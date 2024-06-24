package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pakohan/go/sqlrepo"
)

type modelInfo struct {
	ProjectImportPath string
	Schema            string
	TableName         string
	Columns           []column
}

type column struct {
	Name       string `db:"name"`
	DataType   string `db:"data_type"`
	IsNullable bool   `db:"is_nullable"`
}

func getModelInfo(db *sqlx.DB, q *sqlrepo.SQLRepository, projectImportPath, schema, table string) (*modelInfo, error) {
	columns := []column{}
	err := db.Select(&columns, q.Query("get_columns"), schema, table)
	if err != nil {
		return nil, err
	}

	return &modelInfo{
		ProjectImportPath: projectImportPath,
		Schema:            schema,
		TableName:         table,
		Columns:           columns,
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
