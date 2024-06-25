INSERT INTO "{{.Schema}}"."{{.Model.TableName}}"
({{range $index, $element := (slice .Columns 1)}}
{{- if gt $index 0}}, {{end -}}
 "{{.Name}}"
{{- end}})
VALUE
({{range $index, $element := (slice .Columns 1)}}
{{- if gt $index 0}},  {{end -}}
 :{{.Name}}
{{- end}})
RETURNING {{`{{template "columns.sql"}}`}}