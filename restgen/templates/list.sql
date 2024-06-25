{{`{{template "select.sql"}}`}}
ORDER BY id
{{`{{if .PageSize}}
LIMIT :page_size
{{if .Page}}
OFFSET (:page - 1) * :page_size
{{end}}
{{end}}`}}
