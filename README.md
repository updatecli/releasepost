# releasepost

*Repost changelogs to your static website*

Releasepost can generate json, markdown, and asciidoctor files.

Releasepost creates one file per release note, plus one additional index file.

## Config

releasepost relies on a configuration file to identify what release to monitor and where to repost it.
The configuration file uses the YAML format and accepts the following parameters:

### Changelogs

- `name` defines the changelog name to monitor. It is used to identify a changelog configuration among differentts ones.
- `dir` defines the directory where to mirror the release note
- `kind` defines the kind of changelog to mirror, for example `github`
- `formats` defines a list of format parameters for the mirrored changelog.
- `spec` defines specific parameters to a kind of mirror.


The following example can be used for a HUGO website

**Command**
```
export GITHUB_TOKEN=<insert read-only token>
releasepost --config .releaserepost.yaml
```

**Configuration**

.releaserepost.yaml
```
changelogs:
  - kind: github
    dir: content/en/docs/changelogs/updatecli
    formats:
      - extension: asciidoc
        indexfilename: _index
      - extension: json
        indexfilename: _index
    spec:
      owner: updatecli
      repository: udash
```

### GitHub

#### Credentials

GitHub integration requires a read-only personal access token.
The token must have enough permission to read release information.

#### Configuration

```
changelogs:
  - kind: github
    spec:
      # Define the GitHub owner
      owner: updatecli
      # Define the GitHub repository
      repository: updatecli
      # Define the release type to retrieve
      typefilter:
        draft: false
        prerelease: false
        release: true
        latest: true
      # Define the GitHub url
      url: https://github.com
      # Define the username used to authenticate
      username: john
      # Define the token used to authenticate
      token: xxx
```

The following environment variables will be used as a fallback

* `GITHUB_REPOSITORY` used to set owner and repository
* `GITHUB_TOKEN` used to set the token
* `GITHUB_URL` used to set the GitHub url
* `GITHUB_ACTOR` used to set the GitHub username

## Integration

### Updatecli

Releasepost is designed to work with [Updatecli](https://github.com/updatecli/updatecli) where releasepost is responsible to generate the correct files based on third changelogs and Updatecli to automate the process of publishing them to a git repository.

It can be used with one of the two Updatecli policy:

#### HUGO

Used by the Updatecli project to publish all release notes on [updatecli.io](https://www.updatecli.io/changelogs/updatecli/)

* [ghcr.io/updatecli/policies/releasepost/releasepost](https://github.com/updatecli/policies/tree/main/updatecli/policies/releasepost/releasepost)
* [releasepost.yaml](https://github.com/updatecli/website/blob/master/.releasepost.yaml)
* [GitHub action](https://github.com/updatecli/website/blob/master/.github/workflows/updatecli.weekly.yaml)

#### Docusaurus - rancher.fleet.io

Used by the Rancher Fleet project to publish all release note associated to a MAJOR.MINOR version - [rancher.fleet.io](https://fleet.rancher.io/0.9/changelogs/changelogs/v0.9.3)

* [ghcr.io/olblak/rancherlabs-policies/docusaurus/releasepost](https://github.com/olblak/rancherlabs-policies/tree/main/updatecli/policies/docusaurus/releasepost)
* [releasepost.yaml](https://github.com/rancher/fleet-docs/blob/main/releasepost.yaml)
* [GitHub action](https://github.com/rancher/fleet-docs/blob/main/.github/workflows/updatecli.yml)

## Contributing

This project already solves our initial problem which was to republish Updatecli changelogs from GitHub release to a static website.
That being said, it was also built with extensibility in mind so we could monitor other sources than GitHub releases.

As a community-oriented project, all contributions are greatly appreciated!

Here is a non-exhaustive list of possible contributions:

* ⭐️ this repository.
* Propose a new feature request.
* Highlight an existing feature request with 👍.
* Contribute to any repository in the updatecli organization
* Share the love

## FAQ

**Can releasepost generate changelog?**

No, there are already great tools for doing that.

* [Release Drafter](https://github.com/release-drafter/release-drafter) can automatically generate the next changelogs based on pullrequest labels. If the generated changelog is wrong, you can still update labels on already merged pullrequest and then retrigger release drafter to update the the generated changelog.

* [Conventional Changelog](https://github.com/conventional-changelog/conventional-changelog) can automatically generate the next changelogs based on commit following conventional commit. If the generated changelog is wrong, you can still modify your git history...

* [Changie](https://github.com/miniscruff/changie) is another great tool to generate changelog

The purpose of releasepost is to retrieve already published changelogs and to republish them, for example on a project website.


**Why using _index instead of index on HUGO project?**

[The difference between `index.md` and `_index.md`](https://gohugo.io/content-management/page-bundles/)
