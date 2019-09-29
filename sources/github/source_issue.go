package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/DizzyDoug/extra-extra/connectors/github"
	"github.com/DizzyDoug/extra-extra/sources/types"
)

type issueInfo struct {
	id           int
	closed       bool
	commentCount int
	locked       bool
	lockedReason string
	name         string
}

var (
	_ types.Source = &SourceIssue{} // this is just here to check that SourceIssue implements the Source interface
)

// SourceIssue allows you to track an issue in a repository
// Current trackable items:
// - new comments on specified issue
// - wether issue was locked/unlocked and the reason
// - wether the issue was closed
type SourceIssue struct {
	sourceName   string
	repoOwner    string
	repoName     string
	trackedIssue issueInfo
	connector    *github.Connector
}

// NewSourceIssue is the constructor for a Github Source.
func NewSourceIssue(name, githubLink string,
	issueID int, connector *github.Connector) (SourceIssue, error) {
	repoOwner, repoName, err := parseGithubLink(githubLink)
	if err != nil {
		return SourceIssue{}, err
	}

	source := SourceIssue{
		sourceName: name,
		repoName:   repoName,
		repoOwner:  repoOwner,
		trackedIssue: issueInfo{
			id: issueID,
		},
		connector: connector,
	}

	return source, nil
}

// InitSource initialies the github source, by setting the initial state
func (s *SourceIssue) InitSource() error {
	// initIssues

	return nil
}

// CheckForUpdates check for updates on Issue
func (s *SourceIssue) CheckForUpdates() error {
	return nil
}

// GetName returns the name of this source.
// Source names are unique identifiers to identify a source.
func (s SourceIssue) GetName() string {
	return s.sourceName
}

// GithubApi

type issueResponse struct {
	URL               string        `json:"url"`
	RepositoryURL     string        `json:"repository_url"`
	LabelsURL         string        `json:"labels_url"`
	CommentsURL       string        `json:"comments_url"`
	EventsURL         string        `json:"events_url"`
	HTMLURL           string        `json:"html_url"`
	ID                int           `json:"id"`
	NodeID            string        `json:"node_id"`
	Number            int           `json:"number"`
	Title             string        `json:"title"`
	Labels            []interface{} `json:"labels"`
	State             string        `json:"state"`
	Locked            bool          `json:"locked"`
	Assignee          interface{}   `json:"assignee"`
	Assignees         []interface{} `json:"assignees"`
	Milestone         interface{}   `json:"milestone"`
	Comments          int           `json:"comments"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ClosedAt          time.Time     `json:"closed_at"`
	AuthorAssociation string        `json:"author_association"`
	Body              string        `json:"body"`
}

func (s *SourceIssue) getIssueInfo() (issueResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d", s.repoOwner, s.repoName, s.trackedIssue.id)
	issueResp := issueResponse{}
	res, err := s.connector.Get(url)
	if err != nil {
		return issueResp, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return issueResp, fmt.Errorf("Couldn't get issue info: %d %s", res.StatusCode, res.Status)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return issueResp, err
	}
	err = json.Unmarshal(bz, &issueResp)
	if err != nil {
		return issueResp, err
	}

	return issueResp, nil
}
