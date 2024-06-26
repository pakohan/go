package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/golang-cz/textcase"
	"github.com/jmoiron/sqlx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed templates
	templates embed.FS

	formatSource = true
)

type templateCommand struct {
	filepath     string
	templateName string
	rewrite      bool
}

func generateFiles(cfg *Config) error {
	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"pascal":             textcase.PascalCase,
		"title":              cases.Title(language.Und).String,
		"plural":             plural,
		"remove_underscores": removeUnderscores,
	})
	tmpl, err := tmpl.ParseFS(templates, "templates/*")
	if err != nil {
		return err
	}

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, model := range cfg.Models {
		err = generateModelFiles(db, tmpl, cfg, model)
		if err != nil {
			return err
		}
	}

	return generateGlobalFiles(tmpl, cfg)
}

func generateModelFiles(db *sqlx.DB, tmpl *template.Template, cfg *Config, model Model) error {
	log.Printf("generating model for table %s.%s\n", cfg.TableSchema, model.TableName)

	modelPackage := removeUnderscores(model.TableName)
	commands := []templateCommand{
		{filepath.Join(cfg.BaseDir, "controller", plural(modelPackage), "controller.go"), "controller.go_", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "model.go"), "model.go_", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "columns.sql"), "columns.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "delete.sql"), "delete.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "get.sql"), "get.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "insert.sql"), "insert.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "list.sql"), "list.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "select.sql"), "select.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "update.sql"), "update.sql", true},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "count.sql"), "count.sql", true},
	}

	model.ProjectImportPath = cfg.ProjectImportPath
	mi, err := getModelInfo(db, repo, cfg.TableSchema, model)
	if err != nil {
		return fmt.Errorf("err getting model info: %s", err.Error())
	}

	if len(mi.Columns) == 0 {
		return fmt.Errorf("schema did not return any columns")
	} else if pkColumn := mi.Columns[0]; pkColumn.Name != "id" {
		return fmt.Errorf("first column must be primary key with name id")
	}

	for _, tc := range commands {
		err = executeTemplate(tc, tmpl, mi)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateGlobalFiles(tmpl *template.Template, cfg *Config) error {
	commands := []templateCommand{
		{filepath.Join(cfg.BaseDir, "cmd", "api", "main.go"), "main.go_", false},

		{filepath.Join(cfg.BaseDir, "controller", "controller.go"), "main_controller.go_", true},
		{filepath.Join(cfg.BaseDir, "model", "model.go"), "main_model.go_", true},

		{filepath.Join(cfg.BaseDir, "service", "service.go"), "main_service.go_", true},
		{filepath.Join(cfg.BaseDir, "service", "example", "service.go"), "example_service.go_", true},
	}

	for _, tc := range commands {
		err := executeTemplate(tc, tmpl, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeTemplate(tc templateCommand, tmpl *template.Template, templateParam interface{}) error {
	exists, err := checkFileExists(tc.filepath)
	if err != nil {
		return err
	} else if exists && !tc.rewrite {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(tc.filepath), fs.ModePerm)
	if err != nil {
		return fmt.Errorf("err making directory in file '%s': %s", tc.filepath, err.Error())
	}

	buf := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buf, tc.templateName, templateParam)
	if err != nil {
		return fmt.Errorf("err executing template in file '%s': %s", tc.filepath, err.Error())
	}

	data := buf.Bytes()
	data2 := data
	if filepath.Ext(tc.filepath) == ".go" && formatSource {
		data2, err = format.Source(data)
		if err != nil {
			fmt.Println(string(data))
			return fmt.Errorf("err formatting source in file '%s': %s", tc.filepath, err.Error())
		}
	}

	log.Printf("creating file '%s'\n", tc.filepath)
	err = os.WriteFile(tc.filepath, data2, 0644)
	if err != nil {
		return fmt.Errorf("err writing file in file '%s': %s", tc.filepath, err.Error())
	}

	return nil
}
