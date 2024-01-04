package changelog

import (
	"bytes"
	"fmt"
	"text/template"
)

var (
	// defaultIndexFileName is the default front matters using yaml syntax to add to the index file.
	defaultIndexFrontMatters = ""

	// defaultFrontMatters is the default front matters using yaml syntax to add to the changelog file.
	defaultFrontMatters = ""
)

func renderFrontMatters(data interface{}, frontMatters string) (string, error) {

	tmpl, err := template.New("frontmatters").Parse(frontMatters)
	if err != nil {
		return "", fmt.Errorf("parsing front matters template: %v", err)
	}

	b := bytes.Buffer{}

	if err := tmpl.Execute(&b, data); err != nil {
		return "", fmt.Errorf("executing template: %v", err)
	}

	return b.String(), nil
}
