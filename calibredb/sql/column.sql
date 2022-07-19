{{define "BookCol"}}

{{- $lib := . -}}

IFNULL(JSON_QUOTE(lower(id)), '""') id,
IFNULL(JSON_QUOTE(lower(series_index)), '""') position,
{{range $field := .Fields.BookCol -}}
	IFNULL(JSON_QUOTE({{$field}}), '""') {{GetJsonField $field}},
{{end -}}

{{- range $field := .Fields.DateCol -}}
	IFNULL(JSON_QUOTE(strftime('%Y-%m-%d', {{$field}})), '""') {{GetJsonField $field}},
{{end -}}

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

JSON_QUOTE(
title || IFNULL(" [" || (
	SELECT name 
	FROM series 
	WHERE series.id 
	IN (
		SELECT series 
		FROM books_series_link 
		WHERE book=books.id
	)
) || ", Book " || (series_index) || "]", "")) titleAndSeries, 

IFNULL(JSON_QUOTE((
SELECT text 
FROM comments 
WHERE book=books.id)), '""') description,

IFNULL((
SELECT JSON_QUOTE(lower(rating))
FROM ratings 
WHERE ratings.id 
IN (
	SELECT rating 
	FROM books_ratings_link 
	WHERE book=books.id
)), '""') rating,

{{end}}
