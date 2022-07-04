{{define "categoryField"}}
{{- $label := .Label -}}
{{- $f := getJsonField .Label -}}

{{- if eq $f "description" -}}

IFNULL(JSON_QUOTE((
SELECT text 
FROM comments 
WHERE book=books.id)
), "") {{$f}},

{{- else if eq $f "rating" -}}

JSON_QUOTE(IFNULL((
SELECT lower(rating) 
FROM ratings 
WHERE ratings.id 
IN (
	SELECT rating 
	FROM books_ratings_link 
	WHERE book=books.id
)), "")) {{$f}},

{{- else -}}

{{- if ne .Table "data" -}}
IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(

{{- range $col := .TableColumns -}}
	'value', {{$col}}, 
{{- end -}}

{{- if eq .Table "series" -}}
	'position', lower(series_index),
{{- end -}}

{{- if eq .Table "rating" -}}
	'rating', lower(rating)
{{- end -}}

	'uri', "{{$label}}/" || id,
	'id', lower(id)
{{- end -}}

	))
FROM {{.Table}} 
{{if ne .LinkColumn "" -}}
WHERE {{.Table}}.id 
IN (
	SELECT {{.LinkColumn}}
	FROM books_{{.Table}}_link 
{{end -}}

WHERE book=books.id

{{- if ne .LinkColumn "" -}} 
	) 
{{- end -}}

), "[]") {{$f}}, 

{{- end}}
{{end}}
