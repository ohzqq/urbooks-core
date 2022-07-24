{{define "fieldMeta"}}
'{{.}}', JSON_OBJECT(
'IsMultiple',
CASE JSON_EXTRACT(val, "$.{{.}}.is_multiple")
WHEN '{}' THEN JSON("false")
ELSE JSON("true")
END,
'IsNames', 
CASE JSON_EXTRACT(val, "$.{{.}}.is_multiple.ui_to_list")
WHEN "&" THEN JSON("true")
ELSE JSON("false")
END,
'IsCustom', 
CASE JSON_EXtract(val, "$.{{.}}.is_custom")
WHEN true then JSON("true")
ELSE JSON("false")
END,
'IsEditable', 
CASE JSON_EXtract(val, "$.{{.}}.is_editable")
WHEN true then JSON("true")
ELSE JSON("false")
END,
'IsCategory', 
CASE JSON_EXtract(val, "$.{{.}}.is_category")
WHEN true then JSON("true")
ELSE JSON("false")
END)
{{end}}
