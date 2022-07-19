{{define "CustCol"}}
JSON_OBJECT(
{{range $col := .CustCols -}}

{{- if ne $col.join_table "" -}}

"{{$col.label}}",
IFNULL((
SELECT 
JSON_OBJECT(
'data',
JSON_GROUP_ARRAY(
	JSON_OBJECT(
		'value', value, 
		'id', lower({{$col.table}}.id), 
		'uri', "{{$col.label}}/" || {{$col.table}}.id
	)
),
'meta', JSON_OBJECT(
'is_names',
CASE IFNULL(JSON_EXTRACT(display, "$.is_names"), 0)
WHEN 0 THEN "false"
WHEN 1 THEN "true"
END,
'is_multiple',
CASE is_multiple
WHEN true THEN "true"
ELSE "false"
END))
	
FROM {{$col.table}}, custom_columns 
WHERE {{$col.table}}.id 
IN (
	SELECT value
	FROM {{$col.join_table}}
	WHERE book=books.id
)), '[]'),

{{- else -}}

"{{$col.label}}",
JSON_OBJECT(
'data',
IFNULL((
SELECT 
JSON_QUOTE(value)
FROM {{$col.table}}
WHERE book=books.id
), ''),
'meta', 
(
	SELECT JSON_OBJECT(
'is_names',
CASE IFNULL(JSON_EXTRACT(display, "$.is_names"), 0)
WHEN 0 THEN "false"
WHEN 1 THEN "true"
END,
'is_multiple',
CASE is_multiple
WHEN true THEN "true"
ELSE "false"
END
FROM custom_columns
WHERE ;w

)
)

{{- end -}}

{{end}}
) customColumns,
{{end}}
