package changelog

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	// asciidocTemplate is the template used to generate asciidoctor changelog files
	asciidocTemplate = `{{ .FrontMatters }}
// Disclaimer: this file is generated, do not edit it manually.

{{ if .Changelog.Author }}
__{{ .Changelog.Author }} released this {{ .Changelog.PublishedAt }} - {{ .Changelog.Tag }}__
{{ end }}

=== Description

---
{{ if .Changelog.DescriptionHTML }}
++++

{{ .Changelog.DescriptionHTML }}

++++
{{ else if .Changelog.Description}}
{{ .Changelog.Description }}
{{ end}}
---

{{ if .Changelog.Assets }}

=== Download

[cols="3,1,1" options="header" frame="all" grid="rows"]
|===
| Name | Created At | Updated At
{{ range $asset := .Changelog.Assets }}
| link:{{ $asset.DownloadURL }}[{{ $asset.Name }}] | {{ $asset.CreatedAt }} | {{ $asset.UpdatedAt }}
{{ end }}
|===

{{ end}}
---
{{ if .Changelog.URL }}
__Information retrieved from link:{{ .Changelog.URL }}[here]__
{{ end}}
`

	indexAsciidocTemplate = `{{ .FrontMatters }}
// Disclaimer: this file is generated, do not edit it manually.
[cols="1,1,1" options="header" frame="ends" grid="rows"]
|===
| Name | Author | Published Time
{{ range $pos, $release := .Changelogs }}
| link:changelogs/{{ $release.Tag }}[{{ $release.Name}}{{ if (eq $pos 0) }}(latest){{ end}}] | {{ $release.Author }} | {{ $release.PublishedAt }}
{{ end }}
|===
`
)

func toAsciidocFile(data ReleaseData, filename string, fileTemplate string) error {

	if fileTemplate == "" {
		fileTemplate = asciidocTemplate
	}

	tmpl, err := template.New("asciidoc").Parse(fileTemplate)
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

func toIndexAsciidocFile(data IndexData, filename string, fileTemplate string) error {

	if fileTemplate == "" {
		fileTemplate = indexAsciidocTemplate
	}

	tmpl, err := template.New("asciidoc").
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
		return fmt.Errorf("creating index markdown file %s: %v", filename, err)
	}

	return nil
}
