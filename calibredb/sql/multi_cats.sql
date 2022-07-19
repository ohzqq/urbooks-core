{{define "MultiCats" -}}

{{- $lib := . -}}

{{- range $field := .Fields.MultiCats -}}
{{- $table := GetFieldMeta $lib $field "table"}}
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(
{{- range $key, $col := GetTableColumns $field $lib.Name}}
	'{{$key}}', {{$col}}, 
{{- end -}}
'id', lower(id)))
FROM {{$table}}
WHERE book=books.id), "[]") {{GetJsonField $field}}, 

{{- end -}}
{{- end -}}

