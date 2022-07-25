{{define "rangeFieldMeta"}}
SELECT
JSON_OBJECT(
{{range . -}}
"{{.}}", JSON_EXTRACT(val, "$." || "{{.}}") 
{{end -}}
) fieldMeta
FROM preferences 
WHERE key = "field_metadata"
{{end}}
