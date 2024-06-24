package main

import (
	"github.com/pakohan/go/modelhelper"
	"github.com/pakohan/go/sqlrepo"
)

type modelInfo struct {
	// TableName is the lower case singular name of the table
	Schema    string
	TableName string
	Columns   []column
}

type column struct {
	Name       string `db:"name"`
	DataType   string `db:"data_type"`
	IsNullable bool   `db:"is_nullable"`
}

func getModelInfo(db modelhelper.DB, q *sqlrepo.SQLRepository, schema, table string) (*modelInfo, error) {
	columns := []column{}
	err := db.Select(&columns, q.Query("get_columns"), schema, table)
	if err != nil {
		return nil, err
	}

	return &modelInfo{
		Schema:    schema,
		TableName: table,
		Columns:   columns,
	}, nil
}
