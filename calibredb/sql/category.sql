{{define "category"}}
{{- $lib := . -}}
{{- $field := $lib.GetField $lib.Request.CatLabel -}}

SELECT 
{{if eq $field.LinkColumn "" -}}
	JSON_GROUP_ARRAY(book) books, 
{{- end -}}

{{if $field.IsCustom -}}
	IFNULL(JSON_QUOTE(value), "") value,
{{- else -}}

{{- range $col := $field.TableColumns -}}
	IFNULL(JSON_QUOTE({{$col}}), "") value,
{{- end}}

{{- end -}}

{{if ne $field.LinkColumn "" -}}
(SELECT
JSON_QUOTE(lower(COUNT(book)))
FROM books_{{$field.Table}}_link
WHERE {{$field.LinkColumn}}={{$field.Table}}.id) books,
{{end -}}

JSON_QUOTE("{{$field.Label}}/" || id) uri,
IFNULL(JSON_QUOTE(lower(id)), "") id
FROM {{$field.Table}}

{{if eq $field.LinkColumn "" -}}
	GROUP BY {{$field.Column}}
{{- end -}}
{{end}}
