// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantOwnerName := mockedPullRequest.Base.Repo.Owner.GetLogin()
	gotOwnerName := utils.GetPullRequestBaseOwnerName(mockedPullRequest)

	assert.Equal(t, wantOwnerName, gotOwnerName)
}

func TestGetPullRequestBaseRepoName(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	mockedPullRequest := mockedEnv.GetPullRequest()
	wantRepoName := mockedPullRequest.Base.Repo.GetName()
	gotRepoName := utils.GetPullRequestBaseRepoName(mockedPullRequest)

	assert.Equal(t, wantRepoName, gotRepoName)
}

func TestGetPullRequestNumber(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListCommentsRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantComments := []*github.IssueComment{
		{Body: github.String("Lorem Ipsum")},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				wantComments,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantFiles := []*github.CommitFile{
		{
			Filename: github.String("default-mock-repo/file1.ts"),
			Patch:    nil,
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(wantFiles))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListReviewersRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantReviewers := &github.Reviewers{
		Users: []*github.User{
			{Login: github.String("mary")},
		},
		Teams: []*github.Team{
			{Slug: github.String("reviewpad-team")},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				wantReviewers,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListCollaboratorsRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantCollaborators := []*github.User{
		{Login: github.String("mary")},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposCollaboratorsByOwnerByRepo,
				wantCollaborators,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListAssigneesRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantAssignees := []*github.User{
		{Login: github.String("jane")},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposAssigneesByOwnerByRepo,
				wantAssignees,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListCommitsRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				wantCommits,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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

func TestGetPullRequestReviews(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	wantReviews := []*github.PullRequestReview{
		{
			State: github.String("COMMENTED"),
		},
		{
			State: github.String("COMMENTED"),
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				wantReviews,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotReviews, err := utils.GetPullRequestReviews(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantReviews, gotReviews)
}

func TestGetPullRequestReviews_WhenRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListPullRequestReviewsFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	mockedPullRequest := mockedEnv.GetPullRequest()
	gotReviews, err := utils.GetPullRequestReviews(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		mockedPullRequest.GetUser().GetLogin(),
		mockedPullRequest.GetBase().GetRepo().GetName(),
		mockedPullRequest.GetNumber(),
	)

	assert.Nil(t, gotReviews)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetPullRequests(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ownerName := "testOrg"
	repoName := "testRepo"

	wantPullRequests := []*github.PullRequest{}

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsByOwnerByRepo,
				wantPullRequests,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	gotReviews, err := utils.GetPullRequests(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		ownerName,
		repoName,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantPullRequests, gotReviews)
}

func TestGetPullRequests_WhenRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListPullRequests"

	ownerName := "testOrg"
	repoName := "testRepo"

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepo,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	gotReviews, err := utils.GetPullRequests(
		mockedEnv.GetCtx(),
		mockedEnv.GetClient(),
		ownerName,
		repoName,
	)

	assert.Nil(t, gotReviews)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestGetReviewThreads_WhenRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "GetReviewThreads"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, failMessage, http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	gotThreads, err := utils.GetReviewThreads(
		mockedEnv.GetCtx(),
		mockedEnv.GetClientGQL(),
		aladino.DefaultMockPrOwner,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrNum,
		2,
	)

	assert.Nil(t, gotThreads)
	assert.Equal(t, err.Error(), fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestGetReviewThreads(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedGraphQLQuery := fmt.Sprintf(
		"{\"query\":\"query($pullRequestNumber:Int!$repositoryName:String!$repositoryOwner:String!$reviewThreadsCursor:String){repository(owner: $repositoryOwner, name: $repositoryName){pullRequest(number: $pullRequestNumber){reviewThreads(first: 10, after: $reviewThreadsCursor){nodes{isResolved,isOutdated},pageInfo{endCursor,hasNextPage}}}}}\",\"variables\":{\"pullRequestNumber\":%d,\"repositoryName\":\"%s\",\"repositoryOwner\":\"%s\",\"reviewThreadsCursor\":null}}\n",
		aladino.DefaultMockPrNum,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrOwner,
	)

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			query := aladino.MustRead(req.Body)
			switch query {
			case mockedGraphQLQuery:
				aladino.MustWrite(
					w,
					`{"data": {
                        "repository": {
                            "pullRequest": {
                                "reviewThreads": {
                                    "nodes": [{
                                        "isResolved": true,
                                        "isOutdated": false
                                    }]
                                }
                            }
                        }
                    }}`,
				)
			}
		},
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	wantReviewThreads := []utils.GQLReviewThread{{
		IsResolved: true,
		IsOutdated: false,
	}}
	gotReviewThreads, err := utils.GetReviewThreads(
		mockedEnv.GetCtx(),
		mockedEnv.GetClientGQL(),
		aladino.DefaultMockPrOwner,
		aladino.DefaultMockPrRepoName,
		aladino.DefaultMockPrNum,
		2,
	)

	assert.Nil(t, err)
	assert.Equal(t, gotReviewThreads, wantReviewThreads)

}
