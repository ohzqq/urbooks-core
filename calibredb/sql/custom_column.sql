{{define "customColumn"}}
IFNULL((
SELECT 
{{if .is_multiple -}}

JSON_GROUP_ARRAY(JSON_OBJECT('value', value, 'id', lower(id), 'uri', "{{.Label}}/" || id))
	
FROM {{.table}} 
WHERE {{.table}}.id 
IN (
	SELECT value
	FROM books_{{.table}}_link 
	WHERE book=books.id
)), '[]') {{.label}},

{{- else -}}

JSON_QUOTE(value)
FROM {{.table}}
WHERE book=books.id
), '""') {{.label}},

{{- end -}}
{{end}}
