package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"errors"
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
	"github.com/pakohan/go/sqlrepo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed sql
	sqlDir embed.FS
	repo   = sqlrepo.New(log.Default(), sqlDir, "sql")
)

const configFile = "restgen.json"

type Config struct {
	BaseDir           string `json:"base_dir"`
	ProjectImportPath string `json:"project_import_path"`

	DatabaseURL string  `json:"database_url"`
	TableSchema string  `json:"table_schema"`
	Models      []Model `json:"models"`
}

type Model struct {
	ProjectImportPath string `json:"-"`
	TableName         string `json:"table_name" db:"table_name"`
}

func main() {
	path, err := getConfigPath()
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := readConfig(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = generateFiles(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("go mod init %s && go mod tidy", cfg.ProjectImportPath)
}

func getConfigPath() (string, error) {
	fmt.Printf("checking for config file at '%s'\n", configFile)
	exists, err := checkFileExists(configFile)
	if err != nil {
		return "", err
	} else if exists {
		return configFile, nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Please enter folder for your project: ")
	scanner.Scan()
	baseDir, err := filepath.Abs(scanner.Text())
	if err != nil {
		return "", err
	}

	p := filepath.Join(baseDir, configFile)
	fmt.Printf("checking for config file at '%s'\n", p)
	exists, err = checkFileExists(p)
	if err != nil {
		return "", err
	} else if exists {
		return p, nil
	}

	cfg := &Config{
		BaseDir: baseDir,
	}

	fmt.Print("Please enter postgres connection URL: ")
	scanner.Scan()
	cfg.DatabaseURL = scanner.Text()

	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return "", err
	}
	defer db.Close()

	fmt.Print("Please enter schema: ")
	scanner.Scan()
	cfg.TableSchema = strings.TrimSpace(scanner.Text())

	cfg.Models, err = getTables(db, repo, cfg.TableSchema)
	if err != nil {
		return "", err
	}

	fmt.Print("Please enter base import path for project: ")
	scanner.Scan()
	cfg.ProjectImportPath = strings.TrimSpace(scanner.Text())

	p = filepath.Join(baseDir, configFile)
	log.Printf("Creating json config file '%s'\n", p)
	err = os.MkdirAll(baseDir, fs.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return p, json.NewEncoder(f).Encode(cfg)
}

func readConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
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

type templateCommand struct {
	filepath     string
	templateName string
	rewrite      bool
}

func generateModelFiles(db *sqlx.DB, tmpl *template.Template, cfg *Config, model Model) error {
	log.Printf("generating model for table %s.%s\n", cfg.TableSchema, model.TableName)

	model.ProjectImportPath = cfg.ProjectImportPath
	modelPackage := removeUnderscores(model.TableName)
	commands := []templateCommand{
		{filepath.Join(cfg.BaseDir, "controller", plural(modelPackage), "controller.go"), "controller.go_", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "model.go"), "model.go_", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "columns.sql"), "columns.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "delete.sql"), "delete.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "get.sql"), "get.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "insert.sql"), "insert.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "list.sql"), "list.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "select.sql"), "select.sql", false},
		{filepath.Join(cfg.BaseDir, "model", modelPackage, "sql", "update.sql"), "update.sql", false},
	}

	mi, err := getModelInfo(db, repo, cfg.ProjectImportPath, cfg.TableSchema, model.TableName)
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

		{filepath.Join(cfg.BaseDir, "controller", "controller.go"), "main_controller.go_", false},
		{filepath.Join(cfg.BaseDir, "model", "model.go"), "main_model.go_", false},

		{filepath.Join(cfg.BaseDir, "service", "service.go"), "main_service.go_", false},
		{filepath.Join(cfg.BaseDir, "service", "example", "service.go"), "example_service.go_", false},
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
	if filepath.Ext(tc.filepath) == ".go" {
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

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("err checking if file exists '%s': %s", path, err.Error())
	}
}

func removeUnderscores(s string) string {
	return strings.ReplaceAll(s, "_", "")
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
