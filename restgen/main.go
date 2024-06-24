package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-cz/textcase"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pakohan/go/modelhelper"
	"github.com/pakohan/go/sqlrepo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed sql
	sqlDir embed.FS
)

type Config struct {
	DatabaseURL string
	BaseDir     string
	Models      []Model
}

type Model struct {
	Schema    string
	TableName string
}

func main() {
	cfg := Config{
		DatabaseURL: "postgresql://postgres:test@localhost/postgres?sslmode=disable",
		BaseDir:     "/Users/mogli/Code/go-api",
		Models: []Model{
			{
				Schema:    "main",
				TableName: "example",
			},
		},
	}

	tmpl := template.New("")
	tmpl = tmpl.Funcs(template.FuncMap{
		"pascal":             textcase.PascalCase,
		"title":              cases.Title(language.Und).String,
		"plural":             plural,
		"remove_underscores": removeUnderscores,
	})
	tmpl, err := tmpl.ParseFS(templates, "templates/*")
	if err != nil {
		log.Fatalf("err parsing template: %s\n", err.Error())
	}

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("err opening database connection: %s\n", err.Error())
	}
	defer db.Close()
	mh := modelhelper.DB{DB: db}

	for _, model := range cfg.Models {
		err = generateFiles(mh, tmpl, cfg.BaseDir, model)
		if err != nil {
			log.Fatalf("err generating model files: %s\n", err.Error())
		}
	}
}

func generateFiles(mh modelhelper.DB, tmpl *template.Template, baseDir string, model Model) error {
	modelPackage := removeUnderscores(model.TableName)
	paths := map[string]string{
		filepath.Join(baseDir, "controller", plural(modelPackage), "controller.go"): "controller.go_",
		filepath.Join(baseDir, "model", modelPackage, "model.go"):                   "model.go_",
		filepath.Join(baseDir, "model", modelPackage, "sql", "columns.sql"):         "columns.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "delete.sql"):          "delete.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "get.sql"):             "get.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "insert.sql"):          "insert.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "list.sql"):            "list.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "select.sql"):          "select.sql",
		filepath.Join(baseDir, "model", modelPackage, "sql", "update.sql"):          "update.sql",
	}

	mi, err := getModelInfo(mh, sqlrepo.New(log.Default(), sqlDir, "sql"), model.Schema, model.TableName)
	if err != nil {
		return fmt.Errorf("err getting model info: %s", err.Error())
	}

	if len(mi.Columns) == 0 {
		return fmt.Errorf("schema did not return any columns")
	} else if pkColumn := mi.Columns[0]; pkColumn.Name != "id" {
		return fmt.Errorf("first column must be primary key with name id")
	}

	for path, templateName := range paths {
		err := os.MkdirAll(filepath.Dir(path), fs.ModePerm)
		if err != nil {
			return fmt.Errorf("err making directory in file '%s': %s", path, err.Error())
		}

		buf := &bytes.Buffer{}
		err = tmpl.ExecuteTemplate(buf, templateName, mi)
		if err != nil {
			return fmt.Errorf("err executing template in file '%s': %s", path, err.Error())
		}

		data := buf.Bytes()
		if filepath.Ext(path) == ".go" {
			data, err = format.Source(data)
			if err != nil {
				return fmt.Errorf("err formatting source in file '%s': %s", path, err.Error())
			}
		}

		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return fmt.Errorf("err writing file in file '%s': %s", path, err.Error())
		}
	}
	return nil
}

func removeUnderscores(s string) string {
	return strings.ReplaceAll(s, "_", "-")
}

func plural(s string) string {
	if s == "" {
		return s
	}

	if s[len(s)-1] == 'y' && len(s) > 1 {
		switch s[len(s)-2] {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		// do nothing, since 'y' only gets replaced by 'ie' for plurals in the English language,
		// if the preceding character is not a vowel
		default:
			s = s[:len(s)-1] + "ie"
		}
	}

	return s + "s"
}
