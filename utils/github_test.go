// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/stretchr/testify/assert"
)

type result struct {
	pageNum int
}

func TestGetPullRequestOwnerName(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantOwnerName := mockedPullRequest.Base.Repo.Owner.GetLogin()
	gotOwnerName := utils.GetPullRequestOwnerName(mockedPullRequest)

	assert.Equal(t, wantOwnerName, gotOwnerName)
}

func TestGetPullRequestRepoName(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantRepoName := mockedPullRequest.Base.Repo.GetName()
	gotRepoName := utils.GetPullRequestRepoName(mockedPullRequest)

	assert.Equal(t, wantRepoName, gotRepoName)
}

func TestGetPullRequestNumber(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
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
		return result{}
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
		return result{
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

			return result{pageNum: 1}, resp, nil
		}
		
		return nil, nil, errors.New(failMessage)
	}

	res, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, res)
	assert.EqualError(t, err, failMessage)
}

func TestPaginatedRequest(t *testing.T) {
	initFn := func() interface{} {
		return []*result{
			{pageNum: 1},
		}
	}
	reqFn := func(i interface{}, page int) (interface{}, *github.Response, error) {
		results := i.([]*result)
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

	wantRes := []*result{{pageNum: 1}}
	gotRes, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, err)
	assert.Equal(t, gotRes, wantRes)
}

func TestParseNumPagesFromLink_WhenURLHasNoInfo(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=1>"

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenParseURLFails(t *testing.T) {
	link := "<invalid%+url>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenPageQueryIsNotPresentInURL(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink_WhenNumPagesParseIntFails(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=7B316>; rel=\"last\""

	wantNumPages := 0

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPagesFromLink(t *testing.T) {
	link := "<https://api.github.com/user/58276/repos?page=3>; rel=\"last\""

	// The number of pages are given in the url query "page"
	wantNumPages := 3

	gotNumPages := utils.ParseNumPagesFromLink(link)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestParseNumPages_WhenNoLinkIsProvided(t *testing.T) {
	respHeader := make(http.Header)
	respHeader.Add("Link", " ")
	resp := &github.Response{
		Response: &http.Response{
			Header: respHeader,
		},
	}

	// The number of pages are given in the url query "page"
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

	// The number of pages are given in the url query "page"
	wantNumPages := 3

	gotNumPages := utils.ParseNumPages(resp)

	assert.Equal(t, wantNumPages, gotNumPages)
}

func TestGetPullRequestComments_WhenListCommentsRequestFails(t *testing.T) {
	failMessage := "ListCommentsRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	mockedPullRequestOwner := mockedPullRequest.Base.Repo.Owner.GetLogin()
	mockedPullRequestRepoName := mockedPullRequest.Base.Repo.GetName()
	mockedPullRequestNumber := mockedPullRequest.GetNumber()

	comments, err := utils.GetPullRequestComments(
		mockedEnv.GetCtx(), 
		mockedEnv.GetClient(), 
		mockedPullRequestOwner, 
		mockedPullRequestRepoName, 
		mockedPullRequestNumber, 
		&github.IssueListCommentsOptions{},
	)

	assert.Nil(t, comments)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetPullRequestComments(t *testing.T) {
	wantComments := []*github.IssueComment{
		{Body: github.String("Lorem Ipsum")},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatch(
			mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
			wantComments,
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	mockedPullRequestOwner := mockedPullRequest.Base.Repo.Owner.GetLogin()
	mockedPullRequestRepoName := mockedPullRequest.Base.Repo.GetName()
	mockedPullRequestNumber := mockedPullRequest.GetNumber()

	gotComments, err := utils.GetPullRequestComments(
		mockedEnv.GetCtx(), 
		mockedEnv.GetClient(), 
		mockedPullRequestOwner, 
		mockedPullRequestRepoName, 
		mockedPullRequestNumber, 
		&github.IssueListCommentsOptions{},
	)

	assert.Nil(t, err)
	assert.Equal(t, wantComments, gotComments)
}

func TestGetPullRequestFiles(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	mockedPullRequestOwner := mockedPullRequest.Base.Repo.Owner.GetLogin()
	mockedPullRequestRepoName := mockedPullRequest.Base.Repo.GetName()
	mockedPullRequestNumber := mockedPullRequest.GetNumber()

	wantFiles := []*github.CommitFile{
		{
			Filename: github.String(fmt.Sprintf("%v/file1.ts", mockedPullRequestRepoName)),
			Patch:    nil,
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file2.ts", mockedPullRequestRepoName)),
			Patch:    nil,
		},
		{
			Filename: github.String(fmt.Sprintf("%v/file3.ts", mockedPullRequestRepoName)),
			Patch:    nil,
		},
	}

	gotFiles, err := utils.GetPullRequestFiles(
		mockedEnv.GetCtx(), 
		mockedEnv.GetClient(), 
		mockedPullRequestOwner, 
		mockedPullRequestRepoName, 
		mockedPullRequestNumber,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantFiles, gotFiles)
}
