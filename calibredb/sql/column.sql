{{define "column"}}
{{- $f := GetJsonField .Label -}}
{{- if ne $f "id"}}
{{- if ne $f "url"}}
{{- if ne $f "cover" -}}

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

IFNULL(JSON_QUOTE(

{{- if eq $f "modified" "added" "published" -}}
	strftime('%Y-%m-%d', {{- .Label -}})
{{- else if eq $f "uri" -}}
	"books/" || id
{{- else if eq $f "position" -}}
	lower({{.Label}})
{{- else -}}
	{{- .Label -}}
{{- end -}}

), '""') {{$f}},

{{- end -}}
{{- end -}}
{{- end -}}
{{end}}

