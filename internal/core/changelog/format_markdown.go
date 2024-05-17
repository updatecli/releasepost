package changelog

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	// markdownTemplate is the template used to generate markdown files
	markdownTemplate = `{{ .FrontMatters }}
{{ if .Changelog.Author }}
*{{ .Changelog.Author }} released this {{ .Changelog.PublishedAt }}*
{{ end }}


## Description

{{ if .Changelog.DescriptionHTML }}
{{ .Changelog.DescriptionHTML }}
{{ else if .Changelog.Description}}
{{ .Changelog.Description }}
{{ end}}

{{ if .Changelog.Assets }}
## Download

{{ range $asset := .Changelog.Assets }}
* [{{ $asset.Name }}]({{ $asset.DownloadURL }})
{{ end }}
{{ end}}

{{ if .Changelog.URL }}
*Information retrieved from [here]({{.Changelog.URL}})*
{{ end}}
`

	indexMarkownTemplate = `{{ .FrontMatters }}
{{ range $pos, $release := .Changelogs }}
* [{{ $release.Name}}](changelogs/{{ $release.Tag }}) {{ if (eq $pos 0) }}(latest){{ end}}
{{ end }}
`
)

func toMarkdownFile(data ReleaseData, filename string, fileTemplate string) error {

	if fileTemplate == "" {
		fileTemplate = markdownTemplate
	}

	tmpl, err := template.New("markdown").
		Funcs(sprig.FuncMap()).
		Parse(fileTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	b := bytes.Buffer{}

	if err := tmpl.Execute(&b, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	err = dataToFile(b.Bytes(), filename)
	if err != nil {
		return fmt.Errorf("creating markdown file %s: %v", filename, err)
	}

	return nil
}

func toIndexMarkdownFile(data IndexData, filename string, fileTemplate string) error {

	if fileTemplate == "" {
		fileTemplate = indexMarkownTemplate
	}

	tmpl, err := template.New("markdown").Parse(fileTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}

	b := bytes.Buffer{}
	if err := tmpl.Execute(&b, data); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	err = dataToFile(b.Bytes(), filename)
	if err != nil {
		return fmt.Errorf("creating index markdown file %s: %v", filename, err)
	}

	return nil
}
