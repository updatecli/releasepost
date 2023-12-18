package github

import (
	"fmt"
	"strings"
)

/*
Spec represents the configuration input
*/
type Spec struct {
	/*
		owner specifies the name of a GitHub user or organization.
	*/
	Owner string `yaml:",omitempty" jsonschema:"required"`
	/*
		repository specifies the name of a repository for a specific owner.
	*/
	Repository string `yaml:",omitempty" jsonschema:"required"`
	/*
		TypeFilter specifies the GitHub Release type to retrieve.
	*/
	TypeFilter *ReleaseType `yaml:",omitempty"`
	/*
		"token" specifies the credential used to authenticate with GitHub API.
	*/
	Token string `yaml:",omitempty" jsonschema:"required"`
	/*
		url specifies the default github url in case of GitHub enterprise

		default:
			github.com
	*/
	URL string `yaml:",omitempty"`
	/*
		"username" specifies the username used to authenticate with GitHub API.

		remark:
			the token is usually enough to authenticate with GitHub API.
	*/
	Username string `yaml:",omitempty"`
}

/*
validate verifies if mandatory GitHub parameters are provided and return false otherwise.
*/
func (s *Spec) validate() (errs []error) {
	required := []string{}

	if len(s.Owner) == 0 {
		required = append(required, "owner")
	}

	if len(s.Repository) == 0 {
		required = append(required, "repository")
	}

	if len(required) > 0 {
		errs = append(errs, fmt.Errorf("github parameter(s) required: [%v]", strings.Join(required, ",")))
	}

	return errs
}
