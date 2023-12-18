package changelog

import "time"

// Spec contains various information used to describe target changes
type Spec struct {
	Assets          []AssetSpec `json:",omitempty"`
	Author          string      `json:",omitempty"`
	Title           string      `json:",omitempty"`
	Description     string      `json:",omitempty"`
	DescriptionHTML string      `json:",omitempty"`
	URL             string      `json:",omitempty"`
	PublishedAt     string      `json:",omitempty"`
	UpdatedAt       string      `json:",omitempty"`
	Tag             string      `json:",omitempty"`
	Name            string      `json:",omitempty"`
	/*
		Path is the path to the file containing the changelog.
	*/
	Path string `json:",omitempty"`
}

// ReleaseAssetspec contains information about release assets
type AssetSpec struct {
	CreatedAt   time.Time `json:",omitempty"`
	ContentType string    `json:",omitempty"`
	DownloadURL string    `json:",omitempty"`
	Name        string    `json:",omitempty"`
	UpdatedAt   string    `json:",omitempty"`
	Size        int       `json:",omitempty"`
}
