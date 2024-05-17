package changelog

var (
	// defaultFormats is the default format used to generate changelog files
	defaultFormats = []ConfigFormat{
		{
			Extension:     "markdown",
			IndexFileName: "index",
		},
	}
)

type ConfigFormat struct {
	//
	// Extension is the file extension used for the changelog file.
	//
	// accepted values:
	//    * markdown
	//    * json
	//    * asciidoc
	//
	//  default: markdown
	//
	Extension string
	// FileTemplate is the template used to generate the changelog file
	//
	// default: depends on the extension
	//
	// remark: This setting is useless for json files
	//
	FileTemplate string
	// indexFileName is the name of the index file name without the extension.
	//
	// default: '_index'
	//
	IndexFileName string
	// IndexFileTemplate is the template used to generate the index file
	//
	// default: depends on the extension
	//
	// remark: This setting is useless for json files
	//
	IndexFileTemplate string
	// indexFrontMatters is the front matters using yaml syntax to add to the index file.
	//
	// default: empty
	//
	IndexFrontMatters string
	// frontmatters is the front matters using yaml syntax to add to the changelog file.
	//
	// default: empty
	//
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
