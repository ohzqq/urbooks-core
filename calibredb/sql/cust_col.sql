{{define "custCol"}}
JSON_OBJECT(
{{range $col := .CustCols}}
{{if ne $col.JoinTable "" -}}
"{{$col.Label}}",
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT('value', value, 'id', lower(id), 'uri', "{{$col.Label}}/" || id))
	
FROM {{$col.Table}} 
WHERE {{$col.Table}}.id 
IN (
	SELECT value
	FROM {{$col.JoinTable}}
	WHERE book=books.id
)), '[]'),

{{- else -}}

"{{$col.Label}}",
IFNULL((
SELECT 
JSON_QUOTE(value)
FROM {{$col.Table}}
WHERE book=books.id
), '') 

{{- end -}}

{{end}}
) customColumns,
{{end}}
