{{define "Prefs" -}}
{{- $lib := . -}}
SELECT
JSON_GROUP_OBJECT(
key, 
CASE key
WHEN 'field_metadata' THEN
JSON_OBJECT(
	{{$lib.RenderFieldMetaSql}}
)
ELSE JSON(val)
END
) as data
FROM preferences 
WHERE key 
IN (
	'saved_searches', 
	'field_metadata', 
	'book_display_fields', 
	'tag_browser_hidden_categories'
)
{{end}}
