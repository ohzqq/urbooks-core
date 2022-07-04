{{define "book"}}

{{- $lib := . -}}

SELECT
JSON_QUOTE("{{.Name}}") AS library,

{{- range $f := .Request.Fields -}}
	{{- $field := $lib.GetField $f}}


	{{- if $field.IsCustom -}}
		{{- template "customColumn" $field}}
	{{- else -}}
		{{- if eq $field.Table "" -}}

{{- if eq $field.Label "cover" -}}
JSON_ARRAY(
CASE has_cover
WHEN true
THEN JSON_OBJECT(
	'basename', 'cover',
	'extension', 'jpg',
	'path', "{{$lib.Path}}" || "/" || path || "/cover.jpg",
	'uri', "books/" || books.id || "/cover.jpg",
	'value', 'cover.jpg'
)
ELSE JSON_QUOTE('{}')
END) cover, 
{{end -}}

			{{- template "column" $field}}
		{{- else -}}

{{- if eq $field.Table "data" -}}
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(
	'basename', name,
	'extension', lower(format),
	'path', "{{$lib.Path}}" || "/" || path || "/" || name || '.' || lower(format),
	'size', lower(uncompressed_size),
	'uri', "books/" || books.id,
	'value', name || '.' || lower(format)
{{- end -}}

			{{- template "categoryField" $field}}
		{{- end}}
	{{- end}}
{{end -}}

IFNULL(JSON_QUOTE(lower(id)), "") id

FROM books
{{end}}
