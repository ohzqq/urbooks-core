{{define "book"}}

{{- $lib := . -}}

SELECT
{{range $f := .Request.Fields -}}
	{{- $field := $lib.GetField $f}}

		{{- if $field.IsCustom -}}
			{{- template "customColumn" $field}}
		{{- else -}}
			{{- if eq $field.Table "" -}}
				{{- if eq $field.Label "cover" -}}
IFNULL(
CASE has_cover
WHEN true
THEN JSON_OBJECT(
	'basename', 'cover',
	'extension', 'jpg',
	'path', "{{$lib.Path}}" || "/" || path || "/cover.jpg",
	'uri', "books/" || books.id || "/cover.jpg",
	'value', 'cover.jpg'
)
END, '{}') cover, 
				{{end -}}
				{{- template "column" $field}}
			{{- else -}}
				{{- template "categoryField" $field}}
			{{- end}}
		{{- end}}
{{end -}}

IFNULL(JSON_QUOTE(lower(id)), '""') id

FROM books
{{end}}
