package github

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
)

type releaseAssetsQuery struct {
	RateLimit  RateLimit
	Repository struct {
		Release queriedReleaseAsset `graphql:"release(tagName: $tagName)"`
	} `graphql:"repository(owner: $owner, name: $repository)"`
}

type releaseAssetNode struct {
	CreatedAt   time.Time
	ContentType string
	DownloadURL string
	Name        string
	UpdatedAt   string
	Size        int
}

type releaseAssetsConnection struct {
	Edges      []releaseAssetEdge
	Nodes      []releaseAssetNode
	TotalCount int
	PageInfo   PageInfo
}

type releaseAssetEdge struct {
	Cursor string
	Node   releaseAssetNode
}

type queriedReleaseAsset struct {
	ReleaseAssets releaseAssetsConnection `graphql:"releaseAssets(last: 2, before: $before)"`
	Name          string
	Url           string
}

/*
getReleaseAssets returns a list of assets for a given release
*/
func (g *Github) getReleaseAssets(versionName string) ([]releaseAssetNode, error) {

	results := []releaseAssetNode{}

	fmt.Printf("Looking for release assets related to version %s from %s/%s\n", versionName, g.Spec.Owner, g.Spec.Repository)
	var query releaseAssetsQuery

	variables := map[string]interface{}{
		"owner":      githubv4.String(g.Spec.Owner),
		"repository": githubv4.String(g.Spec.Repository),
		"tagName":    githubv4.String(versionName),
		"before":     (*githubv4.String)(nil),
	}

	for {
		err := g.client.Query(context.Background(), &query, variables)
		if err != nil {
			fmt.Printf("\t %s", err)
			return nil, err
		}

		query.RateLimit.Show()

		for i := len(query.Repository.Release.ReleaseAssets.Edges) - 1; i >= 0; i-- {
			node := query.Repository.Release.ReleaseAssets.Edges[i]

			results = append(results, node.Node)
		}

		if !query.Repository.Release.ReleaseAssets.PageInfo.HasPreviousPage {
			break
		}
		variables["before"] = githubv4.NewString(githubv4.String(query.Repository.Release.ReleaseAssets.PageInfo.StartCursor))
	}

	fmt.Printf("\t=> found %d assets\n", len(results))

	return results, nil

}
