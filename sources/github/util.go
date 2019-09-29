package github

import (
	"errors"
	"strings"
)

// parseGithubLink returns repoOwner, repoName and error
func parseGithubLink(link string) (string, string, error) {
	// removing whitespace
	link = strings.TrimSpace(link)
	// check if link starts with http or https and remove prefix if it exists
	if strings.HasPrefix(link, "http://") {
		link = strings.TrimPrefix(link, "http://")
	} else if strings.HasPrefix(link, "https://") {
		link = strings.TrimPrefix(link, "https://")
	}

	if strings.HasPrefix(link, "github.com") && strings.Count(link, "/") == 2 {
		linkElements := strings.Split(link, "/")
		if len(linkElements) != 3 {
			return "", "", errors.New("Invalid Github link")
		}
		return linkElements[1], linkElements[2], nil
	}

	return "", "", errors.New("Invalid Github link")
}
