{{define "column"}}
{{- $f := GetJsonField .Label -}}
{{- if ne $f "id"}}
{{- if ne $f "url"}}
{{- if ne $f "cover" -}}

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

