UPDATE "{{.Schema}}"."{{.Model.TableName}}"
SET {{range $index, $element := (slice .Columns 1)}}
{{- if gt $index 0}}, {{end -}}
 "{{.Name}}" = :{{.Name}}
{{end -}}
WHERE id = :id
RETURNING {{`{{template "columns.sql"}}`}}