{{define "book"}}
SELECT
{{template "BookCol" . -}}
{{template "JoinCats" . -}}
{{template "SingleCats" . -}}
{{template "MultiCats" . -}}
{{template "CustCol" .}}
JSON_QUOTE("{{.Name}}") library
FROM books
{{end}}
