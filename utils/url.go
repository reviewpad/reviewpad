package utils

import (
	"errors"
	"regexp"

	pbe "github.com/reviewpad/api/go/entities"
)

func ValidateUrl(fileUrl string) (*pbe.Branch, string, error) {
	re := regexp.MustCompile(`^https:\/\/github\.com\/([^/]+)/([^/]+)/blob/([^/]+)/(.+)$`)
	result := re.FindStringSubmatch(fileUrl)
	if len(result) != 5 {
		return nil, "", errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go")
	}

	branch := &pbe.Branch{
		Repo: &pbe.Repository{
			Owner: result[1],
			Name:  result[2],
		},
		Name: result[3],
	}

	return branch, result[4], nil
}
