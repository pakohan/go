package sqlrepo

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"text/template"
)

// SQLRepository
type SQLRepository struct {
	l    *log.Logger
	tmpl *template.Template
}

// New inits a new SQLRepository by calling template.ParseFS on the passed file system.
// Parameter fs is expected to be pointing to a directory with the name passed in dir.
// TODO: check whether the dir param can be dropped by analyzing the fs.FS before calling template.ParseFS.
func New(l *log.Logger, fs fs.FS, dir string) *SQLRepository {
	tmpl, err := template.ParseFS(fs, fmt.Sprintf("%s/*.sql", dir))
	if err != nil {
		l.Fatalf("err parsing template: %s\n", err.Error())
	}

	return &SQLRepository{
		l:    l,
		tmpl: tmpl,
	}
}

// Query executes the sql file as a template and returns the resulting string.
// Parameter param can either be zero or one param, all following will be ignored.
func (q *SQLRepository) Query(name string, param ...interface{}) string {
	if filepath.Ext(name) == "" {
		name += ".sql"
	}

	var templateParam interface{}
	switch len(param) {
	default:
		q.l.Printf("please pass either zero or one param to SQLRepository.Query(). Received %d, will ignore the last %d", len(param), len(param)-1)
		fallthrough // we still want to use the first param.
	case 1:
		templateParam = param[0]
	case 0: // template will be executed with a nil param
	}

	buf := &bytes.Buffer{}
	err := q.tmpl.ExecuteTemplate(buf, name, templateParam)
	if err != nil {
		q.l.Printf("err executing template '%s': %s\n", name, err.Error())
	}

	return buf.String()
}
