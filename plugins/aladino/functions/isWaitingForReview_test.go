// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	host "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var isWaitingForReview = plugins_aladino.PluginBuiltIns().Functions["isWaitingForReview"].Code

func TestIsWaitingForReview_WhenListCommitsRequestFails(t *testing.T) {
	failMessage := "ListCommitsRequestFail"
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
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
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			RequestedReviewers: &pbc.RequestedReviewers{},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	gotValue, gotErr := isWaitingForReview(mockedEnv, args)

	assert.Nil(t, gotValue)
	assert.Equal(t, failMessage, gotErr.(*github.ErrorResponse).Message)
}

func TestIsWaitingForReview_WhenGetPullRequestLastPushDateRequestFails(t *testing.T) {
	failMessage := "GetPullRequestLastPushDateRequestFail"

	mockedLastCommitDate := time.Now()
	mockedCommits := []*github.RepositoryCommit{{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &github.Timestamp{Time: mockedLastCommitDate},
			},
		},
	}}

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedCommits))
				}),
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, failMessage, http.StatusNotFound)
		},
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			RequestedReviewers: &pbc.RequestedReviewers{},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	gotValue, gotErr := isWaitingForReview(mockedEnv, args)

	assert.Nil(t, gotValue)
	assert.EqualError(t, gotErr, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestIsWaitingForReview_WhenListReviewsRequestFails(t *testing.T) {
	failMessage := "ListReviewsRequestFail"

	mockedLastCommitDate := time.Now()
	mockedCommits := []*github.RepositoryCommit{{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &github.Timestamp{Time: mockedLastCommitDate},
			},
		},
	}}

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		RequestedReviewers: &pbc.RequestedReviewers{},
	})

	mockedCodeReviewNumber := host.GetPullRequestNumber(mockedCodeReview)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)

	mockedPullRequestLastPushDateGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int! $repositoryName:String! $repositoryOwner:String!){
			repository(owner: $repositoryOwner, name: $repositoryName){
				pullRequest(number: $pullRequestNumber){
					timelineItems(last: 1, itemTypes: [HEAD_REF_FORCE_PUSHED_EVENT, PULL_REQUEST_COMMIT]){
						nodes{
							__typename,
							...on HeadRefForcePushedEvent {
								createdAt
							},
							...on PullRequestCommit {
								commit {
									pushedDate,
									committedDate
								}
							}
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber": %d,
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedCodeReviewNumber, mockRepo, mockOwner)

	mockedPullRequestLastPushDateGQLQueryBody := fmt.Sprintf(`{
		"data": {
			"repository": {
				"pullRequest": {
					"timelineItems": {
						"nodes": [{
							"__typename": "PullRequestCommit",
							"commit": {
								"pushedDate": "%s"
							}
						}]
					}
				}
			}
		}
	}`, mockedLastCommitDate.Format("2006-01-02T15:04:05Z07:00"))

	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedCommits))
				}),
			),
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
		func(w http.ResponseWriter, req *http.Request) {
			query := utils.MinifyQuery(utils.MustRead(req.Body))
			switch query {
			case utils.MinifyQuery(mockedPullRequestLastPushDateGQLQuery):
				utils.MustWrite(w, mockedPullRequestLastPushDateGQLQueryBody)
			}
		},
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			RequestedReviewers: &pbc.RequestedReviewers{},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []lang.Value{}
	gotValue, gotErr := isWaitingForReview(mockedEnv, args)

	assert.Nil(t, gotValue)
	assert.Equal(t, failMessage, gotErr.(*github.ErrorResponse).Message)
}

func TestIsWaitingForReview_WhenNoCommits(t *testing.T) {
	mockedCommits := []*github.RepositoryCommit{}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedCommits))
				}),
			),
		},
		nil,
		aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
			RequestedReviewers: &pbc.RequestedReviewers{},
		}),
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)
	args := []lang.Value{}
	gotValue, err := isWaitingForReview(mockedEnv, args)

	wantValue := lang.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestIsWaitingForReview_WhenHasNoReviews(t *testing.T) {
	mockedLastCommitDate := time.Now()
	tests := map[string]struct {
		requestedReviewers []*pbc.User
		requestedTeams     []*pbc.Team
		wantValue          lang.Value
	}{
		"from user": {
			requestedReviewers: []*pbc.User{{
				Login: "john",
			}},
			requestedTeams: []*pbc.Team{},
			wantValue:      lang.BuildBoolValue(true),
		},
		"from team": {
			requestedReviewers: []*pbc.User{},
			requestedTeams: []*pbc.Team{{
				Name: "team",
			}},
			wantValue: lang.BuildBoolValue(true),
		},
		"from user and team": {
			requestedReviewers: []*pbc.User{{
				Login: "john",
			}},
			requestedTeams: []*pbc.Team{{
				Name: "team",
			}},
			wantValue: lang.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(
								[]*github.RepositoryCommit{{
									Commit: &github.Commit{
										Committer: &github.CommitAuthor{
											Date: &github.Timestamp{Time: mockedLastCommitDate},
										},
									},
								}},
							))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal([]*github.PullRequestReview{}))
						}),
					),
				},
				nil,
				aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
					RequestedReviewers: &pbc.RequestedReviewers{
						Users: test.requestedReviewers,
						Teams: test.requestedTeams,
					},
				}),
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			)
			args := []lang.Value{}
			gotValue, err := isWaitingForReview(mockedEnv, args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantValue, gotValue)
		})
	}
}

func TestIsWaitingForReview_WhenHasReviews(t *testing.T) {
	mockedAuthorLogin := "author"
	mockedLastCommitDate := time.Now()
	mockedCreateDateAfterLastCommitDate := mockedLastCommitDate.Add(time.Hour)
	mockedCreateDateBeforeLastCommitDate := mockedLastCommitDate.Add(time.Hour * -1)
	mockedCommits := []*github.RepositoryCommit{{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &github.Timestamp{Time: mockedLastCommitDate},
			},
		},
	}}
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		RequestedReviewers: &pbc.RequestedReviewers{},
		Author:             &pbc.User{Login: mockedAuthorLogin},
	})

	mockedCodeReviewNumber := host.GetPullRequestNumber(mockedCodeReview)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedCodeReview)
	mockRepo := host.GetPullRequestBaseRepoName(mockedCodeReview)
	mockedPullRequestLastPushDateGQLQuery := fmt.Sprintf(`{
		"query":"query($pullRequestNumber:Int! $repositoryName:String! $repositoryOwner:String!){
			repository(owner: $repositoryOwner, name: $repositoryName){
				pullRequest(number: $pullRequestNumber){
					timelineItems(last: 1, itemTypes: [HEAD_REF_FORCE_PUSHED_EVENT, PULL_REQUEST_COMMIT]){
						nodes{
							__typename,
							...on HeadRefForcePushedEvent {
								createdAt
							},
							...on PullRequestCommit {
								commit {
									pushedDate,
									committedDate
								}
							}
						}
					}
				}
			}
		}",
		"variables":{
			"pullRequestNumber": %d,
			"repositoryName": "%s",
			"repositoryOwner": "%s"
		}
	}`, mockedCodeReviewNumber, mockRepo, mockOwner)

	mockedPullRequestLastPushDateGQLQueryBody := fmt.Sprintf(`{
		"data": {
			"repository": {
				"pullRequest": {
					"timelineItems": {
						"nodes": [{
							"__typename": "PullRequestCommit",
							"commit": {
								"pushedDate": "%s"
							}
						}]
					}
				}
			}
		}
	}`, mockedLastCommitDate.Format("2006-01-02T15:04:05Z07:00"))

	tests := map[string]struct {
		reviewers []*github.PullRequestReview
		wantValue lang.Value
	}{
		"from nil user": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: nil,
				},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from author": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String(mockedAuthorLogin),
				},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from user but outdated CHANGES_REQUESTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("CHANGES_REQUESTED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateBeforeLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(true),
		},
		"from user but outdated COMMENTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("COMMENTED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateBeforeLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(true),
		},
		"from user but outdated APPROVED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("APPROVED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateBeforeLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from user and up to date CHANGES_REQUESTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("CHANGES_REQUESTED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateAfterLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from user and up to date COMMENTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("COMMENTED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateAfterLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from user and up to date APPROVED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("APPROVED"),
				SubmittedAt: &github.Timestamp{Time: mockedCreateDateAfterLastCommitDate},
			}},
			wantValue: lang.BuildBoolValue(false),
		},
		"from user and more than one": {
			reviewers: []*github.PullRequestReview{
				{
					User: &github.User{
						Login: github.String("user"),
					},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &github.Timestamp{Time: mockedCreateDateBeforeLastCommitDate},
				},
				{
					User: &github.User{
						Login: github.String("user"),
					},
					State:       github.String("APPROVED"),
					SubmittedAt: &github.Timestamp{Time: mockedCreateDateAfterLastCommitDate},
				},
			},
			wantValue: lang.BuildBoolValue(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(mockedCommits))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(test.reviewers))
						}),
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(utils.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedPullRequestLastPushDateGQLQuery):
						utils.MustWrite(w, mockedPullRequestLastPushDateGQLQueryBody)
					}
				},
				mockedCodeReview,
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			)
			args := []lang.Value{}
			gotValue, err := isWaitingForReview(mockedEnv, args)

			assert.Nil(t, err, "Expected no error")
			assert.True(t, gotValue.Equals(test.wantValue), "Expected %v, got %v", test.wantValue, gotValue)
		})
	}
}
