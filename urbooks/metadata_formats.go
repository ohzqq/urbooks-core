package urbooks

import (
	"log"
	"os"
	"text/template"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gosimple/slug"
)

func (b *Book) ToFFmeta() {
	meta, err := os.Create(slug.Make(b.Get("title").String()) + ".ini")
	if err != nil {
		log.Fatal(err)
	}
	defer meta.Close()

	err = MetaFmt.FFmeta.Execute(meta, b)
	if err != nil {
		log.Fatal(err)
	}
}

//func (b *Book) ToPlain() string {
//  var buf bytes.Buffer
//  err := MetaFmt.Plain.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  return buf.String()
//}

//func (b *Book) ToMarkdown() string {
//  var buf bytes.Buffer
//  err := MetaFmt.MD.Execute(&buf, b)
//  if err != nil {
//    log.Fatal(err)
//  }
//  //fmt.Println(markdown)
//  return buf.String()
//}

type metadataFormats struct {
	FFmeta *template.Template
	MD     *template.Template
	Plain  *template.Template
	Cue    *template.Template
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
	MD:     template.Must(template.New("html").Funcs(funcMap).Parse(mdTmpl)),
	Plain:  template.Must(template.New("plain").Funcs(funcMap).Parse(plainTmpl)),
}

const ffmetaTmpl = `;FFMETADATA
{{$title := .Get "titleAndSeries" -}}
title={{$title.String}}
album={{$title.String}}
artist=
{{- with $authors := .Get "authors" -}}
	{{- $authors.String -}}
{{- end}}
composer=
{{- with $narrators := .Get "narrators" -}}
	{{- $narrators.String -}}
{{- end}}
genre=
{{- with $tags := .Get "tags" -}}
	{{- $tags.String -}}
{{- end}}
comment=
{{- with $description := .Get "description" -}}
	{{- $description.String -}}
{{- end -}}
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
