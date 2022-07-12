{{define "categoryField"}}
{{- $label := .Label -}}
{{- $f := GetJsonField .Label -}}

{{- if eq $f "description" -}}

IFNULL(JSON_QUOTE((
SELECT text 
FROM comments 
WHERE book=books.id)
), '""') {{$f}},

{{- else if eq .Table "ratings" -}}

IFNULL((
SELECT JSON_QUOTE(lower(rating))
FROM ratings 
WHERE ratings.id 
IN (
	SELECT rating 
	FROM books_ratings_link 
	WHERE book=books.id
)), '""') {{$f}},

{{- else -}}

	{{- if ne .Table "ratings" -}}

IFNULL((
SELECT 
		{{if .IsMultiple -}}
			JSON_GROUP_ARRAY(
		{{- end -}}

			JSON_OBJECT(

		{{- range $key, $col := .TableColumns -}}
			'{{$key}}', {{$col}}, 
		{{- end -}}

			'id', lower(id)
	{{- end -}}

	{{if .IsMultiple -}}
		)
	{{- end -}}
)
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

	{{- if .IsMultiple -}}
		), "[]") {{$f}}, 
	{{- else -}}
		), "{}") {{$f}}, 
	{{- end -}}

	{{- end}}
{{end}}
