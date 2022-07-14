// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/stretchr/testify/assert"
)

type paginatedRequestResult struct {
    pageNum int
}

func TestGetPullRequestOwnerName(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantOwnerName := mockedPullRequest.Base.Repo.Owner.GetLogin()
	gotOwnerName := utils.GetPullRequestOwnerName(mockedPullRequest)

	assert.Equal(t, wantOwnerName, gotOwnerName)
}

func TestGetPullRequestRepoName(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantRepoName := mockedPullRequest.Base.Repo.GetName()
	gotRepoName := utils.GetPullRequestRepoName(mockedPullRequest)

	assert.Equal(t, wantRepoName, gotRepoName)
}

func TestGetPullRequestNumber(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantPullRequestNumber := mockedPullRequest.GetNumber()
	gotPullRequestNumber := utils.GetPullRequestNumber(mockedPullRequest)

	assert.Equal(t, wantPullRequestNumber, gotPullRequestNumber)
}

func TestPaginatedRequest_WhenFirstRequestFails(t *testing.T) {
	failMessage := "PaginatedRequestFail"
	initFn := func() interface{} {
		return paginatedRequestResult{}
	}
	reqFn := func(i interface{}, page int) (interface{}, *github.Response, error) {
		return nil, nil, errors.New(failMessage)
	}

	res, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, res)
	assert.EqualError(t, err, failMessage)
}

func TestPaginatedRequest_WhenFurtherRequestsFail(t *testing.T) {
	failMessage := "PaginatedRequestFail"
	initFn := func() interface{} {
		return paginatedRequestResult{
			pageNum: 1,
		}
	}
	reqFn := func(i interface{}, page int) (interface{}, *github.Response, error) {
		if page == 1 {
			respHeader := make(http.Header)
			respHeader.Add("Link", "<https://api.github.com/user/58276/repos?page=3>; rel=\"last\"")
			resp := &github.Response{
				Response: &http.Response{
					Header: respHeader,
				},
				NextPage: 3,
			}

			return paginatedRequestResult{pageNum: 1}, resp, nil
		}

		return nil, nil, errors.New(failMessage)
	}

	res, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, res)
	assert.EqualError(t, err, failMessage)
}

func TestPaginatedRequest(t *testing.T) {
	initFn := func() interface{} {
		return []*paginatedRequestResult{
			{pageNum: 1},
		}
	}
	reqFn := func(i interface{}, page int) (interface{}, *github.Response, error) {
		results := i.([]*paginatedRequestResult)
		if page == 1 {
			respHeader := make(http.Header)
			respHeader.Add("Link", "<https://api.github.com/user/58276/repos?page=3>; rel=\"last\"")
			resp := &github.Response{
				Response: &http.Response{
					Header: respHeader,
				},
			}

			return results, resp, nil
		}

		return results, nil, nil
	}

	wantRes := []*paginatedRequestResult{{pageNum: 1}}
	gotRes, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, err)
	assert.Equal(t, gotRes, wantRes)
}

func TestParseNumPagesFromLink_WhenHTTPLinkHeaderHasNoRel(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=1>"

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenHTTPLinkHeaderIsInvalid(t *testing.T) {
	link := "<invalid%+url>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenHTTPLinkHeaderHasNoQueryParamPage(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenHTTPLinkHeaderHasInvalidQueryParamPage(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=7B316>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=3>; rel=\"last\""

	// The number of pages is provided in the url query parameter "page"
	wantNumPages := 3

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPages_WhenHTTPLinkHeaderIsNotProvided(t *testing.T) {
	respHeader := make(http.Header)
	respHeader.Add("Link", " ")
	resp := &github.Response{
		Response: &http.Response{
			Header: respHeader,
		},
	}

	wantNumPages := 0

	gotNumPages := utils.ParseNumPages(resp)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPages(t *testing.T) {
	respHeader := make(http.Header)
	respHeader.Add("Link", "<https://api.github.com/user/58276/repos?page=3>; rel=\"last\"")
	resp := &github.Response{
		Response: &http.Response{
			Header: respHeader,
		},
	}

	// The number of pages is provided in the url query parameter "page"
	wantNumPages := 3

	gotNumPages := utils.ParseNumPages(resp)

	assert.Equal(t, wantNumPages, gotNumPages)
}
