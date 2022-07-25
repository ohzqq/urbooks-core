{{define "CustCol"}}
JSON_OBJECT(
{{range $col := .CustCols}}
{{/* ltrim("{{$col.label}}", '#'), JSON_OBJECT( */}}
"{{$col.label}}", JSON_OBJECT(
'meta', IFNULL(
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
		END)
	FROM custom_columns
	WHERE custom_columns.id = {{$col.id}}
), "{}"),

{{- if ne $col.join_table "" -}}

'data', IFNULL(
(
	SELECT JSON_GROUP_ARRAY(JSON_OBJECT(
		'value', value, 
		'id', lower({{$col.table}}.id), 
		'uri', ltrim("{{$col.label}}/", '#') || {{$col.table}}.id))
	FROM {{$col.table}}, custom_columns 
	WHERE {{$col.table}}.id 
	IN (SELECT value
		FROM {{$col.join_table}}
		WHERE book=books.id)
), '[]')
),

{{- else}}
'data', IFNULL(
(
	SELECT JSON_QUOTE(value)
	FROM {{$col.table}}
	WHERE book=books.id
), '')
)

{{- end -}}
{{end}}
) customColumns,

{{end}}
