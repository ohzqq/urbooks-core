{{define "booksInCategory"}}

{{- $lib := . -}}
{{- $field := $lib.GetField $lib.Request.CatLabel}}

SELECT
{{$field.Column}},
(
	SELECT
	GROUP_CONCAT(book)
	FROM books_{{$field.Table}}_link
	WHERE {{$field.LinkColumn}}={{$field.Table}}.id
) itemIDs
FROM {{$field.Table}} 
WHERE id={{$lib.Request.PathID}}
{{end}}
