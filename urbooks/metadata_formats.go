package urbooks

import (
	"text/template"
	"log"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type metadataFormats struct {
	FFmeta *template.Template
	MD *template.Template
	Plain *template.Template
	Cue *template.Template
}

var funcMap = template.FuncMap{
	"toMarkdown": toMarkdown,
}

func toMarkdown(str string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(str)
	if err != nil {
		log.Fatal(err)
	}
	return markdown
}

var MetaFmt = metadataFormats{
	FFmeta: template.Must(template.New("ffmeta").Parse(ffmetaTmpl)),
	MD: template.Must(template.New("html").Funcs(funcMap).Parse(mdTmpl)),
	Plain: template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
}

const ffmetaTmpl = `;FFMETADATA
title={{if .TitleAndSeries}}{{.TitleAndSeries}}{{end}}
album={{if .TitleAndSeries}}{{.TitleAndSeries}}{{end}}
artist={{if .Authors}}{{.Authors.Join}}{{end}}
composer={{if .Narrators}}{{.Narrators.Join}}{{end}}
genre={{if .Tags}}{{.Tags.Join}}{{end}}
comment={{if .Description}}{{.Description}}{{end}}
`

const mdTmpl = `{{if .Title}}# {{.Title}}   
{{end}}{{if .HasSeries}}**Series:** {{.SeriesString}}   
{{end}}{{if .Authors}}**Authors:** {{.Authors.Join}}   
{{end}}{{if .Narrators}}**Narrators:** {{.Narrators.Join}}   
{{end}}{{if .Tags}}**Tags:** {{.Tags.Join}}   
{{end}}{{if .Rating}}**Rating:** {{.Rating}}   
{{end}}{{if .Description}}**Description:** {{toMarkdown .Description}}{{end}}`

const plainTmpl = `{{if .Title}}{{.Title}}   
{{end}}{{if .HasSeries}}Series: {{.SeriesString}}   
{{end}}{{if .Authors}}Authors: {{.Authors.Join}}   
{{end}}{{if .Narrators}}Narrators: {{.Narrators.Join}}   
{{end}}{{if .Tags}}Tags: {{.Tags.Join}}   
{{end}}{{if .Rating}}Rating: {{.Rating}}   
{{end}}{{if .Description}}Description: {{.Description}}{{end}}`
