{{`{{template "select.sql"}}`}}
WHERE TRUE
{{range $index, $element := .Columns}}
{{- if .IsFilter -}}
{{`{{if`}} .{{pascal .Name}}{{`}}`}}
AND "{{.Name}}" = :{{.Name}}
{{`{{end}}`}}
{{end -}}
{{end}}
ORDER BY id
{{`{{if .PageSize}}
LIMIT :page_size
{{if .Page}}
OFFSET (:page - 1) * :page_size
{{end}}
{{end}}`}}
