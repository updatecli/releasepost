package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/shurcooL/githubv4"
	"github.com/updatecli/updatecli/pkg/plugins/utils/version"
	"golang.org/x/oauth2"
)

/*
GitHub contains settings to interact with GitHub
*/
type Github struct {
	// Spec contains inputs coming from updatecli configuration
	Spec          Spec
	client        gitHubClient
	releaseType   ReleaseType
	versionFilter version.Filter
}

/*
gitHubClient must be implemented by any GitHub query client (v4 API)
*/
type gitHubClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

/*
New returns a new valid GitHub object.
*/
func New(s interface{}) (*Github, error) {

	newSpec := Spec{}

	err := mapstructure.Decode(s, &newSpec)
	if err != nil {
		return nil, err
	}

	errs := newSpec.validate()

	if len(errs) > 0 {
		strErrs := []string{}
		for _, err := range errs {
			strErrs = append(strErrs, err.Error())
		}
		return &Github{}, fmt.Errorf(strings.Join(strErrs, "\n"))
	}

	newSpec.mergeFromEnv("RELEASEPOST_GITHUB")
	newSpec.mergeFromEnv("GITHUB")

	if newSpec.URL == "" {
		newSpec.URL = "github.com"
	}

	if !strings.HasPrefix(newSpec.URL, "https://") && !strings.HasPrefix(newSpec.URL, "http://") {
		newSpec.URL = "https://" + newSpec.URL
	}

	newFilters, err := newSpec.VersionFilter.Init()
	if err != nil {
		return &Github{}, fmt.Errorf("initializing version filter: %w", err)
	}

	// Initialize github client
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: newSpec.Token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	g := Github{
		Spec:          newSpec,
		versionFilter: newFilters,
	}

	if newSpec.TypeFilter != nil {
		g.releaseType = *newSpec.TypeFilter
	} else {
		g.releaseType = ReleaseType{
			Release: true,
		}
	}

	if strings.HasSuffix(newSpec.URL, "github.com") {
		g.client = githubv4.NewClient(httpClient)
	} else {
		// For GH enterprise the GraphQL API path is /api/graphql
		// Cf https://docs.github.com/en/enterprise-cloud@latest/graphql/guides/managing-enterprise-accounts#3-setting-up-insomnia-to-use-the-github-graphql-api-with-enterprise-accounts
		graphqlURL, err := url.JoinPath(newSpec.URL, "/api/graphql")
		if err != nil {
			return nil, err
		}
		g.client = githubv4.NewEnterpriseClient(graphqlURL, httpClient)
	}

	return &g, nil
}

/*
mergeFromEnv updates the target receiver with the "non zero-ed" environment variables
*/
func (gs *Spec) mergeFromEnv(envPrefix string) {
	prefix := fmt.Sprintf("%s_", envPrefix)

	github_repository := os.Getenv(fmt.Sprintf("%s_%s", prefix, "REPOSITORY"))

	if github_repository != "" && strings.Contains(github_repository, "/") {
		repositoryArray := strings.Split(github_repository, "/")

		if repositoryArray[0] != "" && gs.Owner == "" {
			gs.Owner = os.Getenv(fmt.Sprintf("%s%s", prefix, "OWNER"))
		}
		if repositoryArray[1] != "" && gs.Repository == "" {
			gs.Repository = os.Getenv(fmt.Sprintf("%s%s", prefix, "REPOSITORY"))
		}
	}

	if os.Getenv(fmt.Sprintf("%s%s", prefix, "TOKEN")) != "" && gs.Token == "" {
		gs.Token = os.Getenv(fmt.Sprintf("%s%s", prefix, "TOKEN"))
	}
	if os.Getenv(fmt.Sprintf("%s%s", prefix, "URL")) != "" && gs.URL == "" {
		gs.URL = os.Getenv(fmt.Sprintf("%s%s", prefix, "URL"))
	}
	if os.Getenv(fmt.Sprintf("%s%s", prefix, "ACTOR")) != "" && gs.URL == "" {
		gs.Username = os.Getenv(fmt.Sprintf("%s%s", prefix, "USERNAME"))
	}
}
