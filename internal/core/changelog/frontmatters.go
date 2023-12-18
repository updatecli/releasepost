package changelog

import (
	"bytes"
	"fmt"
	"text/template"
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
