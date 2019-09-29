package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/DizzyDoug/extra-extra/connectors/github"
	"github.com/DizzyDoug/extra-extra/sources/types"
)

type releaseInfo struct {
	id          int
	name        string
	tagName     string
	releaseDate string
	description string
	url         string
}

var (
	_ types.Source = &SourceRelease{} // this is just here to check that SourceGithub implements the Source interface
)

// SourceRelease allows you to track the releases of a repository on Github:
type SourceRelease struct {
	sourceName    string
	repoOwner     string
	repoName      string
	latestRelease releaseInfo
	connector     *github.Connector
}

// NewSourceRelease is the constructor for a Github Source.
func NewSourceRelease(name, githubLink string, connector *github.Connector) (SourceRelease, error) {
	repoOwner, repoName, err := parseGithubLink(githubLink)
	if err != nil {
		return SourceRelease{}, err
	}

	source := SourceRelease{
		sourceName: name,
		repoName:   repoName,
		repoOwner:  repoOwner,
		connector:  connector,
	}

	return source, nil
}

// InitSource initialies the github source, by setting the initial state
func (s *SourceRelease) InitSource() error {
	lrr, err := s.getLatestReleaseFromGithub()
	if err != nil {
		return err
	}

	s.latestRelease = releaseInfo{
		id:          lrr.ID,
		name:        lrr.Name,
		description: lrr.Body,
		tagName:     lrr.TagName,
		releaseDate: lrr.PublishedAt,
		url:         lrr.URL,
	}

	return nil
}

// CheckForUpdates checks for new releases
func (s *SourceRelease) CheckForUpdates() error {
	return nil
}

// GetName returns the name of this source.
// Source names are unique identifiers to identify a source.
func (s SourceRelease) GetName() string {
	return s.sourceName
}

// GithubApi

// releaseAPIResponse is the type representing the reponse of the latest release endpoint, some unneeded properites are omitted
type releaseAPIResponse struct {
	URL             string `json:"url"`              // "https://api.github.com/repos/octocat/Hello-World/releases/1",
	HTMLURL         string `json:"html_url"`         // "https://github.com/octocat/Hello-World/releases/v1.0.0",
	AssetURL        string `json:"assets_url"`       // "https://api.github.com/repos/octocat/Hello-World/releases/1/assets",
	UploadURL       string `json:"upload_url"`       // "https://uploads.github.com/repos/octocat/Hello-World/releases/1/assets{?name,label}",
	TarballURL      string `json:"tarball_url"`      // "https://api.github.com/repos/octocat/Hello-World/tarball/v1.0.0",
	ZipballURL      string `json:"zipball_url"`      // "https://api.github.com/repos/octocat/Hello-World/zipball/v1.0.0",
	ID              int    `json:"id"`               // 1,
	NodeID          string `json:"node_id"`          // "MDc6UmVsZWFzZTE=",
	TagName         string `json:"tag_name"`         // "v1.0.0",
	TargetCommitish string `json:"target_commitish"` // "master",
	Name            string `json:"name"`             // "v1.0.0",
	Body            string `json:"body"`             // "Description of the release",
	Draft           bool   `json:"draft"`            // false,
	Prerelease      bool   `json:"prerelease"`       // false,
	CreatedAt       string `json:"created_at"`       // "2013-02-27T19:35:32Z",
	PublishedAt     string `json:"published_at"`     // "2013-02-27T19:35:32Z",
}

func (s *SourceRelease) getLatestReleaseFromGithub() (releaseAPIResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", s.repoOwner, s.repoName)
	releaseResp := releaseAPIResponse{}
	res, err := s.connector.Get(url)
	if err != nil {
		return releaseResp, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return releaseResp, fmt.Errorf("Couldn't get latest release: %d %s", res.StatusCode, res.Status)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return releaseResp, err
	}
	err = json.Unmarshal(bz, &releaseResp)
	if err != nil {
		return releaseResp, err
	}

	return releaseResp, nil
}
