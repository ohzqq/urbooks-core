{{define "customColumn"}}
IFNULL((
SELECT 
{{if .IsMultiple -}}

JSON_GROUP_ARRAY(JSON_OBJECT('value', value, 'id', lower(id), 'uri', "{{.Label}}/" || id))
	
FROM {{.Table}} 
WHERE {{.Table}}.id 
IN (
	SELECT value
	FROM books_{{.Table}}_link 
	WHERE book=books.id
)

{{- else -}}

JSON_QUOTE(value)
FROM {{.Table}}
WHERE book=books.id

{{- end -}}

), "[]") {{.Label}},
{{end}}
