package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

// releasesQuery defines a github v4 API query to retrieve a list of releases sorted by reverse order of created time.
/*
https://developer.github.com/v4/explorer/
# Query
query getLatestRelease($owner: String!, $repository: String!){
	rateLimit {
		cost
		remaining
		resetAt
	}
	repository(owner: $owner, name: $repository){
		releases(last:100, before: $before, orderBy:$orderBy){
			totalCount
			pageInfo {
				hasNextPage
				endCursor
			}
			edges {
				node {
							name
					tagName
				isDraft
				isPrerelease
				}
				cursor
			}
		}
	}
}
# Variables
{
	"owner": "updatecli",
	"repository": "updatecli"
}
*/
type releasesQuery struct {
	RateLimit  RateLimit
	Repository struct {
		Releases repositoryRelease `graphql:"releases(last: 100, before: $before, orderBy: $orderBy)"`
	} `graphql:"repository(owner: $owner, name: $repository)"`
}
type releaseNode struct {
	Name         string
	TagName      string
	IsDraft      bool
	IsLatest     bool
	IsPrerelease bool
}
type releaseEdge struct {
	Cursor string
	Node   releaseNode
}
type repositoryRelease struct {
	TotalCount int
	PageInfo   PageInfo
	Edges      []releaseEdge
}

/*
searchReleases return every releases from the github api
ordered by reverse order of created time.
Draft and pre-releases are filtered out.
*/
func (g *Github) searchReleases() (releases []string, err error) {

	fmt.Println("Searching releases")

	var query releasesQuery

	variables := map[string]interface{}{
		"owner":      githubv4.String(g.Spec.Owner),
		"repository": githubv4.String(g.Spec.Repository),
		"before":     (*githubv4.String)(nil),
		"orderBy": githubv4.ReleaseOrder{
			Field:     "CREATED_AT",
			Direction: "ASC",
		},
	}

	for {
		err := g.client.Query(context.Background(), &query, variables)
		if err != nil {
			fmt.Printf("\t%s", err)
			return releases, err
		}

		query.RateLimit.Show()

		for i := len(query.Repository.Releases.Edges) - 1; i >= 0; i-- {
			node := query.Repository.Releases.Edges[i]

			// If releaseType.Latest is set to true, then it means
			// we only care about identifying the latest release
			if g.releaseType.Latest {
				if node.Node.IsLatest {
					releases = append(releases, node.Node.TagName)
					break
				}
				// Check if the next release is of type "latest"
				continue
			}

			if node.Node.IsDraft {
				if g.releaseType.Draft {
					releases = append(releases, node.Node.TagName)
				}
			} else if node.Node.IsPrerelease {
				if g.releaseType.PreRelease {
					releases = append(releases, node.Node.TagName)
				}
			} else {
				if g.releaseType.Release {
					releases = append(releases, node.Node.TagName)
				}
			}
		}

		if !query.Repository.Releases.PageInfo.HasPreviousPage {
			break
		}

		variables["before"] = githubv4.NewString(githubv4.String(query.Repository.Releases.PageInfo.StartCursor))
	}

	fmt.Printf("%d releases found", len(releases))
	return releases, nil
}
