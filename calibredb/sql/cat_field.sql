{{define "JoinCats"}}

{{- $lib := .Lib -}}

{{- range $field := .JoinCats}}

IFNULL((
SELECT 
JSON_GROUP_ARRAY(JSON_OBJECT(

{{- range $key, $col := GetTableColumns $field $lib -}}
	'{{$key}}', {{$col}}, 
{{- end -}}

'id', lower(id)))
FROM {{.}} 
WHERE {{.}}.id 
IN (
SELECT 
	{{- if eq . "tags" "authors"}} name {{end -}}
	{{- if eq . "languages"}} lang_code {{end -}}
	FROM books_{{.}}_link 
	WHERE book=books.id)
), "[]") {{.}}, 

{{- end -}}
{{- end -}}
