package github

import (
	"fmt"

	"github.com/olblak/releasepost/internal/core/changelog"
)

/*
Run retrieves the changelog from GitHub.
*/
func (g Github) Run() ([]changelog.Spec, error) {

	var releases []string
	var err error
	var changelogs []changelog.Spec

	releases, err = g.searchReleases()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve GitHub releases: %v", err)
	}

	for i := range releases {
		fmt.Printf("Release found: %s\n", releases[i])

		var changelog *changelog.Spec

		changelog, err = g.changelog(releases[i])
		if err != nil {
			fmt.Printf("unable to retrieve GitHub changelog: %v\n", err.Error())
			continue
		}

		changelogs = append(changelogs, *changelog)
	}

	return changelogs, nil
}
