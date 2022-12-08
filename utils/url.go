package utils

import (
	"errors"
	"regexp"

	"github.com/google/go-github/v48/github"
)

func ValidateUrl(fileUrl string) (*github.PullRequestBranch, string, error) {
	re := regexp.MustCompile(`^https:\/\/github\.com\/([^/]+)/([^/]+)/blob/([^/]+)/(.+)$`)
	result := re.FindStringSubmatch(fileUrl)
	if len(result) != 5 {
		return nil, "", errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go")
	}

	branch := &github.PullRequestBranch{
		Repo: &github.Repository{
			Owner: &github.User{
				Login: github.String(result[1]),
			},
			Name: github.String(result[2]),
		},
		Ref: github.String(result[3]),
	}

	return branch, result[4], nil
}
