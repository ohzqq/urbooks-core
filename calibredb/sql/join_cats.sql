{{define "JoinCats" -}}

{{- $lib := . -}}

{{- range $field := .Fields.JoinCats -}}
{{- $table := GetFieldMeta $lib $field "table" -}}
{{- $link := GetFieldMeta $lib $field "link_column" -}}
{{- $join := GetFieldMeta $lib $field "join_table"}}
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(
{{- range $key, $col := GetTableColumns $field $lib.Name}}
	'{{$key}}', {{$col}}, 
{{- end -}}

'id', lower(id)))
FROM {{$table}}
WHERE {{$table}}.id 
IN (
	SELECT {{$link}}
	FROM {{$join}}
	WHERE book=books.id)
), "[]") {{GetJsonField $field}}, 

{{- end -}}
{{- end -}}
