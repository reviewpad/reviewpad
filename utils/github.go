// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/tomnomnom/linkheader"
)

func GetPullRequestOwnerName(pullRequest *github.PullRequest) string {
	return pullRequest.Base.Repo.Owner.GetLogin()
}

func GetPullRequestRepoName(pullRequest *github.PullRequest) string {
	return pullRequest.Base.Repo.GetName()
}

func GetPullRequestNumber(pullRequest *github.PullRequest) int {
	return pullRequest.GetNumber()
}

func PaginatedRequest(
	initFn func() interface{},
	reqFn func(interface{}, int) (interface{}, *github.Response, error),
) (interface{}, error) {
	page := 1
	results, resp, err := reqFn(initFn(), page)
	if err != nil {
		return nil, err
	}

	numPages := ParseNumPages(resp)
	page++
	for page <= numPages && resp.NextPage > page {
		results, _, err = reqFn(results, page)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func ParseNumPagesFromLink(link string) int {
	urlInfo := linkheader.Parse(link).FilterByRel("last")
	if len(urlInfo) < 1 {
		return 0
	}

	urlData, err := url.Parse(urlInfo[0].URL)
	if err != nil {
		return 0
	}

	numPagesStr := urlData.Query().Get("page")
	if numPagesStr == "" {
		return 0
	}

	numPages, err := strconv.ParseInt(numPagesStr, 10, 32)
	if err != nil {
		return 0
	}

	return int(numPages)
}

//ParseNumPages Given a link header string representing pagination info, returns total number of pages.
func ParseNumPages(resp *github.Response) int {
	link := resp.Header.Get("Link")
	if strings.Trim(link, " ") == "" {
		return 0
	}

	return ParseNumPagesFromLink(link)
}

func HasLinearHistory(commit *github.Commit) bool {
	switch len(commit.Parents) {
	case 0:
		return true
	case 1:
		return HasLinearHistory(commit.Parents[0])
	default:
		return false
	}
}
