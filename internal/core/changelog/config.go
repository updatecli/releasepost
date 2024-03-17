package changelog

import (
	"fmt"
	"path/filepath"
)

type IndexData struct {
	Latest       Spec
	Changelogs   []Spec
	FrontMatters string
}

type ReleaseData struct {
	Changelog    Spec
	FrontMatters string
}

var (
	defaultChangelogDir = "changelogs"
	// defaultKind is the default kind of changelog to retrieve
	defaultKind = "github"
	// defaultDir is the default directory where changelog files will be repost
	defaultDir = "dist"

	nameMarkdown = "markdown"
	nameJson     = "json"
	nameAsciidoc = "asciidoc"
)

/*
Config contains various information used to configure the way changelogs are retrieved
*/
type Config struct {
	/*
		Name specifies the name of the changelog to retrieve. It is only used for identification purposes.
	*/
	Name string
	/*
		dir specifies the directory where changelog files will be repost.
	*/
	Dir string
	/*
		kind specifies the kind of changelog to retrieve.

		accepted values:
			* github
	*/
	Kind string
	/*
		format specifies the format of the changelog to generate.
	*/
	Formats []ConfigFormat
	/*
		spec specifies the configuration input for a specific kind of changelog.
	*/
	Spec interface{}
}

/*
Sanitize ensures that the configuration is valid and that all required fields are set.
*/
func (c *Config) Sanitize(configFile string) error {

	if c.Name != "" {
		fmt.Printf("Changelog configuration: %q\n", c.Name)
	}

	if c.Dir == "" {
		fmt.Printf("No directory specified, using default directory %v\n", defaultDir)
		c.Dir = defaultDir
	}

	if c.Formats == nil {
		fmt.Printf("No format specified, using default format %v\n", defaultFormats)
		c.Formats = defaultFormats
	}

	if c.Kind == "" {
		fmt.Printf("No kind specified, using default kind %v\n", defaultKind)
		c.Kind = defaultKind
	}

	// Ensure that any relative file path is relative to the compose file
	sanitizePath := func(path string) string {
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(configFile), path)
		}
		return path
	}

	for i := range c.Formats {
		if err := c.Formats[i].Sanitize(); err != nil {
			fmt.Printf("sanitizing format: %v", err)
			continue
		}
	}

	c.Dir = sanitizePath(c.Dir)

	return nil
}

/*
SaveToDisk saves one changelog per file per format to disk
*/
func (c Config) SaveToDisk(changelogs []Spec) error {

	fmt.Println("Generating changelogs")

	err := initDir(filepath.Join(c.Dir, defaultChangelogDir))
	if err != nil {
		return err
	}

	for i := range changelogs {
		data := ReleaseData{
			Changelog: changelogs[i],
		}

		frontMatters, err := renderFrontMatters(data, c.Formats[0].FrontMatters)
		if err != nil {
			fmt.Printf("rendering front matters: %v", err)
			continue
		}

		data.FrontMatters = frontMatters

		for _, format := range c.Formats {
			switch format.Extension {
			case nameAsciidoc:
				filename := filepath.Join(
					c.Dir,
					defaultChangelogDir,
					changelogs[i].Tag+".adoc",
				)
				if err = toAsciidocFile(data, filename); err != nil {
					fmt.Printf("creating asciidoc file %s: %v", filename, err)
					continue
				}

			case nameJson:
				filename := filepath.Join(
					c.Dir,
					defaultChangelogDir,
					changelogs[i].Tag+".json",
				)
				if err = toJsonFile(data, filename); err != nil {
					fmt.Printf("creating json file %s: %v", filename, err)
					continue
				}

			case nameMarkdown:
				filename := filepath.Join(
					c.Dir,
					defaultChangelogDir,
					changelogs[i].Tag+".md",
				)
				if err = toMarkdownFile(data, filename); err != nil {
					fmt.Printf("creating markdown file %s: %v", filename, err)
					continue
				}
			default:
				fmt.Printf("unknown changelog format %q", format)
				continue
			}
		}
	}

	return nil
}

/*
SaveIndexToDisk saves an index file per format to disk
*/
func (c Config) SaveIndexToDisk(changelogs []Spec) error {

	if len(changelogs) == 0 {
		fmt.Printf("No changelog found for %s\n", c.Name)
		return nil
	}

	fmt.Println("Generating index")

	err := initDir(filepath.Join(c.Dir, defaultChangelogDir))
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	data := IndexData{
		Changelogs: changelogs,
	}

	for _, format := range c.Formats {

		data.FrontMatters, err = renderFrontMatters(data, format.IndexFrontMatters)
		if err != nil {
			fmt.Printf("rendering front matters: %v", err)
			continue
		}

		switch format.Extension {
		case nameAsciidoc:
			indexFileName := format.IndexFileName + ".adoc"
			if err := toIndexAsciidocFile(data, filepath.Join(c.Dir, indexFileName)); err != nil {
				fmt.Printf("creating index asciidoc file %s: %v\n", filepath.Join(c.Dir, indexFileName), err)
			}

		case nameJson:
			indexFileName := format.IndexFileName + ".json"

			/*
				The purpose of shortChangelogs is to generate a json file
				that only contains the bare minimum.
			*/

			shortChangelogs := make([]Spec, len(changelogs))
			for j := range changelogs {
				shortChangelogs[j].Author = changelogs[j].Author
				shortChangelogs[j].Tag = changelogs[j].Tag
				shortChangelogs[j].PublishedAt = changelogs[j].PublishedAt
			}
			data.Changelogs = shortChangelogs
			data.Latest = shortChangelogs[0]
			data.FrontMatters = ""

			if err = toJsonFile(data, filepath.Join(c.Dir, indexFileName)); err != nil {
				fmt.Printf("creating index json file %s: %v", filepath.Join(c.Dir, "index.json"), err)
				continue
			}

		case nameMarkdown:
			indexFileName := format.IndexFileName + ".md"
			if err := toIndexMarkdownFile(data, filepath.Join(c.Dir, indexFileName)); err != nil {
				fmt.Printf("creating index markdown file %s: %v", filepath.Join(c.Dir, indexFileName), err)
				continue
			}

		default:
			fmt.Printf("unknown changelog format %q", format)
			continue
		}
	}

	return nil
}
