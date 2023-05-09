// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/tomnomnom/linkheader"
)

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
	for page <= numPages && resp.NextPage > page {
		page++
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

// ParseNumPages Given a link header string representing pagination info, returns total number of pages.
func ParseNumPages(resp *github.Response) int {
	link := resp.Header.Get("Link")
	if strings.Trim(link, " ") == "" {
		return 0
	}

	return ParseNumPagesFromLink(link)
}
