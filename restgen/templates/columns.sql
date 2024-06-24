{{range $index, $element := .Columns}}
{{- if gt $index 0}}, {{end -}}
 "{{.Name}}"
{{end}}