package result

import "fmt"

type Changelog struct {
	Created    []string `json:"created"`
	Modified   []string `json:"modified"`
	UnModified []string `json:"unmodified"`
}

var (
	ChangelogResult Changelog
)

func (c Changelog) String() {
	if len(c.Created) > 0 {
		fmt.Println("Changelogs reposted:")
		for _, v := range c.Created {
			fmt.Printf("\t* %s\n", v)
		}
	}

	if len(c.Modified) > 0 {
		fmt.Println("Changelogs modified:")
		for _, v := range c.Modified {
			fmt.Printf("\t* %s\n", v)
		}
	}

	if len(c.UnModified) > 0 {
		fmt.Println("Changelogs unmodified:")
		for _, v := range c.UnModified {
			fmt.Printf("\t* %s\n", v)
		}
	}
}

func (c Changelog) ExitCode() int {
	if len(c.Created) == 0 && len(c.Modified) == 0 {
		return 1
	}
	return 0
}
