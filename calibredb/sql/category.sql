{{define "category"}}
{{- $lib := . -}}
{{- $field := $lib.Request.CatLabel -}}
{{- $label := GetFieldMeta $lib $lib.Request.CatLabel "label" -}}
{{- $table := GetFieldMeta $lib $lib.Request.CatLabel "table" -}}
{{- $link := GetFieldMeta $lib $lib.Request.CatLabel "link_column" -}}
{{- $join := GetFieldMeta $lib $lib.Request.CatLabel "join_table" -}}
{{- $custom := GetFieldMeta $lib $lib.Request.CatLabel "custom_column" -}}

{{if eq $field "identifiers" "formats"}}
SELECT
JSON_GROUP_ARRAY(book) books,
{{- range $key, $col := GetTableColumns $field $lib.Name -}}
	{{- if ne $key "path" -}}
		{{- if ne $key "uri"}}
IFNULL(JSON_QUOTE({{$col}}), "") {{$key}}, 
		{{- end -}}
	{{- end -}}
{{- end}}
IFNULL(JSON_QUOTE(lower(id)), '""') id
FROM {{$table}}
{{if eq $table "data" -}}
GROUP BY extension
{{- else if eq $table "identifiers" -}}
GROUP BY val
{{end}}

{{- else -}}

SELECT
{{if eq $custom "true" -}}
IFNULL(JSON_QUOTE(value), '""') value,
IFNULL(JSON_QUOTE("{{$label}}/" || id), '""') uri,
{{- end -}}

{{- range $key, $col := GetTableColumns $field $lib.Name}}
	{{- if ne $key "position"}}
IFNULL(JSON_QUOTE({{$col}}), "") {{$key}}, 
	{{- end -}}
{{- end}}

(SELECT
JSON_QUOTE(lower(COUNT(book)))
FROM books_{{$table}}_link
WHERE {{$link}}={{$table}}.id) books,
IFNULL(JSON_QUOTE(lower(id)), "") id
FROM {{$table}}
{{- end -}}
{{end}}
