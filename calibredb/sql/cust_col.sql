{{define "custCol"}}
JSON_OBJECT(
{{range $col := .CustCols}}
{{if ne $col.join_table "" -}}
"{{$col.label}}",
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT('value', value, 'id', lower(id), 'uri', "{{$col.Label}}/" || id))
	
FROM {{$col.table}} 
WHERE {{$col.table}}.id 
IN (
	SELECT value
	FROM {{$col.join_table}}
	WHERE book=books.id
)), '[]'),

{{- else -}}

"{{$col.label}}",
IFNULL((
SELECT 
JSON_QUOTE(value)
FROM {{$col.table}}
WHERE book=books.id
), '') 

{{- end -}}

{{end}}
) customColumns,
{{end}}
