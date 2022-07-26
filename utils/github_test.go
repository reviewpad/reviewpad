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
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

type paginatedRequestResult struct {
	pageNum int
}

func TestGetPullRequestHeadOwnerName(t *testing.T) {
	mockedHeadOwnerName := "reviewpad"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedHeadOwnerName),
				},
			},
		},
	})
	wantOwnerName := mockedPullRequest.Head.Repo.Owner.GetLogin()
	gotOwnerName := utils.GetPullRequestHeadOwnerName(mockedPullRequest)

	assert.Equal(t, wantOwnerName, gotOwnerName)
	assert.Equal(t, mockedHeadOwnerName, gotOwnerName)
}

func TestGetPullRequestHeadRepoName(t *testing.T) {
	mockedHeadRepoName := "mocks-test"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Name: &mockedHeadRepoName,
			},
		},
	})
	wantRepoName := mockedPullRequest.Head.Repo.GetName()
	gotRepoName := utils.GetPullRequestHeadRepoName(mockedPullRequest)

	assert.Equal(t, wantRepoName, gotRepoName)
	assert.Equal(t, mockedHeadRepoName, gotRepoName)
}

func TestGetPullRequestBaseOwnerName(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantOwnerName := mockedPullRequest.Base.Repo.Owner.GetLogin()
	gotOwnerName := utils.GetPullRequestBaseOwnerName(mockedPullRequest)

	assert.Equal(t, wantOwnerName, gotOwnerName)
}

func TestGetPullRequestBaseRepoName(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantRepoName := mockedPullRequest.Base.Repo.GetName()
	gotRepoName := utils.GetPullRequestBaseRepoName(mockedPullRequest)

	assert.Equal(t, wantRepoName, gotRepoName)
}

func TestGetPullRequestNumber(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
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

func TestGetPullRequestComments_WhenListCommentsRequestFails(t *testing.T) {
	failMessage := "ListCommentsRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
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
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	comments, err := utils.GetPullRequestComments(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
		mockedPullRequest.GetNumber(),
		&github.IssueListCommentsOptions{},
	)

	assert.Nil(t, comments)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetPullRequestComments(t *testing.T) {
	wantComments := []*github.IssueComment{
		{Body: github.String("Lorem Ipsum")},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				wantComments,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotComments, err := utils.GetPullRequestComments(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
		mockedPullRequest.GetNumber(),
		&github.IssueListCommentsOptions{},
	)

	assert.Nil(t, err)
	assert.Equal(t, wantComments, gotComments)
}

func TestGetPullRequestFiles(t *testing.T) {
	wantFiles := []*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file1.ts"),
			Patch:    nil,
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(wantFiles))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotFiles, err := utils.GetPullRequestFiles(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
		mockedPullRequest.GetNumber(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantFiles, gotFiles)
}

func TestGetPullRequestReviewers_WhenListReviewersRequestFails(t *testing.T) {
	failMessage := "ListReviewersRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	reviewers, err := utils.GetPullRequestReviewers(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
		mockedPullRequest.GetNumber(),
		nil,
	)

	assert.Nil(t, reviewers)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetPullRequestReviewers(t *testing.T) {
	wantReviewers := &github.Reviewers{
		Users: []*github.User{
			{Login: github.String("mary")},
		},
		Teams: []*github.Team{
			{Slug: github.String("reviewpad-team")},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				wantReviewers,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotReviewers, err := utils.GetPullRequestReviewers(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
		mockedPullRequest.GetNumber(),
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantReviewers, gotReviewers)
}

func TestGetRepoCollaborators_WhenListCollaboratorsRequestFails(t *testing.T) {
	failMessage := "ListCollaboratorsRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCollaboratorsByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	collaborators, err := utils.GetRepoCollaborators(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
	)

	assert.Nil(t, collaborators)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetRepoCollaborators(t *testing.T) {
	wantCollaborators := []*github.User{
		{Login: github.String("mary")},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposCollaboratorsByOwnerByRepo,
				wantCollaborators,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotCollaborators, err := utils.GetRepoCollaborators(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.Base.Repo.Owner.GetLogin(),
		mockedPullRequest.Base.Repo.GetName(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantCollaborators, gotCollaborators)
}

func TestGetIssuesAvailableAssignees_WhenListAssigneesRequestFails(t *testing.T) {
	failMessage := "ListAssigneesRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposAssigneesByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotAssignees, err := utils.GetIssuesAvailableAssignees(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
	)

	assert.Nil(t, gotAssignees)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetIssuesAvailableAssignees(t *testing.T) {
	wantAssignees := []*github.User{
		{Login: github.String("jane")},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposAssigneesByOwnerByRepo,
				wantAssignees,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotAssignees, err := utils.GetIssuesAvailableAssignees(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantAssignees, gotAssignees)
}

func TestGetPullRequestCommits_WhenListCommistsRequestFails(t *testing.T) {
	failMessage := "ListCommitsRequestFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotCommits, err := utils.GetPullRequestCommits(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
	)

	assert.Nil(t, gotCommits)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetPullRequestCommits(t *testing.T) {
	wantCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
			},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				wantCommits,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotCommits, err := utils.GetPullRequestCommits(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantCommits, gotCommits)
}

func TestGetPullRequestReviewsWithCommentedState(t *testing.T) {
	wantReviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
		},
		{
			State: github.String("COMMENTED"),
		},
	}
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				wantReviews,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotReviews, err := utils.GetPullRequestReviews(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
		&github.ListOptions{},
	)

	assert.Nil(t, err)
	assert.Equal(t, wantReviews, gotReviews)
}

func TestGetPullRequestReviewsFails(t *testing.T) {
	failMessage := "ListPullRequestReviewsFail"
	mockedEnv, err := aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotReviews, err := utils.GetPullRequestReviews(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
		&github.ListOptions{},
	)

	assert.Nil(t, gotReviews)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}
