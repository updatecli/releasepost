package changelog

var (
	// defaultFormats is the default format used to generate changelog files
	defaultFormats = []ConfigFormat{
		{
			Extension:     "markdown",
			IndexFileName: "_index",
		},
	}

	// defaultIndexFileName is the default front matters using yaml syntax to add to the index file.
	defaultIndexFrontMatters = `---
title: Changelogs
---`

	defaultFrontMatters = `---
title: "{{ .Changelog.Name }}"
date: {{ .Changelog.PublishedAt }}
---`
)

type ConfigFormat struct {
	/*
		Extension is the file extension used for the changelog file.
		accepted values:
			* markdown
			* json

		default:
			* markdown
	*/
	Extension string
	/*
		indexFileName is the name of the index file name without the extension.

		default:
			* _index
	*/
	IndexFileName string
	/*
		indexFrontMatters is the front matters using yaml syntax to add to the index file.
	*/
	IndexFrontMatters string
	/*
		frontmatters is the front matters using yaml syntax to add to the changelog file.
	*/
	FrontMatters string
}

func (c *ConfigFormat) Sanitize() error {

	if c.Extension == "" {
		c.Extension = "markdown"
	}
	if c.IndexFileName == "" {
		c.IndexFileName = "index"
	}

	if c.IndexFrontMatters == "" {
		c.IndexFrontMatters = defaultIndexFrontMatters
	}

	if c.FrontMatters == "" {
		c.FrontMatters = defaultFrontMatters
	}

	return nil
}
