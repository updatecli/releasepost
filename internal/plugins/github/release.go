package github

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/shurcooL/githubv4"

	"github.com/updatecli/releasepost/internal/core/changelog"
)

/*
	changelogQuery defines a github v4 API query to retrieve the changelog of a given release
	https://docs.github.com/en/graphql/reference/objects#release

	https://developer.github.com/v4/explorer/

# Query

	query getLatestRelease($owner: String!, $repository: String!){
		repository(owner: $owner, name: $repository){
			release(tagName: "v0.17.0"){
				description
				publishedAt
				url
			}
		}
	}

# Variables

	{
		"owner": "updatecli",
		"repository": "updatecli"
	}
*/
type changelogQuery struct {
	Repository struct {
		Release queriedRelease `graphql:"release(tagName: $tagName)"`
	} `graphql:"repository(owner: $owner, name: $repository)"`
}

type queriedRelease struct {
	Author          user
	Description     string
	DescriptionHTML string `graphql:"descriptionHTML"`
	PublishedAt     time.Time
	Name            string
	Tag             tagRef
	Url             string
	UpdatedAt       time.Time
}

type tagRef struct {
	Name string
}

type user struct {
	Name  string
	Login string
}

/*
changelog returns a changelog description based on a release name
*/
func (g *Github) changelog(versionName string) (*changelog.Spec, error) {

	fmt.Printf("Looking for release information related to version %s from %s/%s\n", versionName, g.Spec.Owner, g.Spec.Repository)
	var query changelogQuery

	variables := map[string]interface{}{
		"owner":      githubv4.String(g.Spec.Owner),
		"repository": githubv4.String(g.Spec.Repository),
		"tagName":    githubv4.String(versionName),
	}

	err := g.client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Printf("\t %s", err)
		return nil, err
	}

	URL, err := url.JoinPath(g.Spec.URL, g.Spec.Owner, g.Spec.Repository)

	if err != nil {
		return nil, err
	}

	if len(query.Repository.Release.Url) == 0 {
		return nil, fmt.Errorf("no GitHub Release found for %s on %q", versionName, URL)
	}

	assets, err := g.getReleaseAssets(versionName)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve assets for %s on %q", versionName, URL)
	}

	releaseAssets := make([]changelog.AssetSpec, len(assets))
	for i := range assets {
		releaseAssets[i].CreatedAt = assets[i].CreatedAt
		releaseAssets[i].ContentType = assets[i].ContentType
		releaseAssets[i].DownloadURL = assets[i].DownloadURL
		releaseAssets[i].Name = assets[i].Name
		releaseAssets[i].UpdatedAt = assets[i].UpdatedAt
		releaseAssets[i].Size = assets[i].Size
	}

	return &changelog.Spec{
		Author:          query.Repository.Release.Author.Name + " (" + query.Repository.Release.Author.Login + ")",
		Assets:          releaseAssets,
		Description:     query.Repository.Release.Description,
		DescriptionHTML: query.Repository.Release.DescriptionHTML,
		PublishedAt:     query.Repository.Release.PublishedAt.String(),
		Name:            query.Repository.Release.Name,
		Tag:             query.Repository.Release.Tag.Name,
		UpdatedAt:       query.Repository.Release.UpdatedAt.String(),
		URL:             query.Repository.Release.Url,
	}, nil
}
