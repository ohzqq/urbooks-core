{{define "book"}}

{{- $lib := . -}}

SELECT

{{- template "custCol" $lib}}

{{range $f := .Request.Fields -}}
	{{- $field := $lib.GetField $f}}

		{{- if $field.IsCustom -}}
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

JSON_QUOTE(
title || 
IFNULL(
	" [" || (
		SELECT name 
		FROM series 
		WHERE series.id 
		IN (
			SELECT series 
			FROM books_series_link 
			WHERE book=books.id
		)
	) ||
	", Book " || 
	(series_index) ||
	"]"
	, "") 
)
titleAndSeries, 

JSON_QUOTE("{{$lib.Name}}") library,

IFNULL(JSON_QUOTE(lower(id)), '""') id

FROM books
{{end}}
