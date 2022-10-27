// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"

	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var review = plugins_aladino.PluginBuiltIns().Actions["review"].Code

type ReviewRequestPostBody struct {
	Event string `json:"event"`
	Body  string `json:"body"`
}

func TestReview_WhenGetAuthenticatedUserRequestFails(t *testing.T) {
	failMessage := "GetAuthenticatedUserRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetUser,
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
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("some comment")}
	err := review(mockedEnv, args)

	assert.Equal(t, failMessage, err.(*github.ErrorResponse).Message)
}

func TestReview_WhenReviewMethodIsUnsupported(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetUser,
				&github.User{
					Login: github.String("bot-account"),
				},
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("INVALID"), aladino.BuildStringValue("some comment")}
	err := review(mockedEnv, args)

	assert.EqualError(t, err, "review: unsupported review event INVALID")
}

func TestReview_WhenListReviewsRequestFails(t *testing.T) {
	failMessage := "ListReviewsRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetUser,
				&github.User{
					Login: github.String("bot-account"),
				},
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
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("some comment")}
	err := review(mockedEnv, args)

	assert.Equal(t, failMessage, err.(*github.ErrorResponse).Message)
}

func TestReview_WhenGetPullRequestLastPushDateRequestFails(t *testing.T) {
	failMessage := "GetPullRequestLastPushDateRequestFail"

	pullRequestSubmittedAt := time.Now()

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetUser,
				&github.User{
					Login: github.String("bot-account"),
				},
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{
					{
						ID:   github.Int64(1),
						Body: github.String(""),
						User: &github.User{
							Login: github.String("bot-account"),
						},
						State:       github.String("COMMENTED"),
						SubmittedAt: &pullRequestSubmittedAt,
					},
				},
			),
		},
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, failMessage, http.StatusNotFound)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{aladino.BuildStringValue("APPROVE"), aladino.BuildStringValue("some comment")}
	err := review(mockedEnv, args)

	assert.EqualError(t, err, fmt.Sprintf("non-200 OK status code: 404 Not Found body: \"%s\\n\"", failMessage))
}

func TestReview_WhenBodyIsNotEmpty(t *testing.T) {
	reviewBody := "test comment"

	mockedLastCommitDate := time.Now()

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestLastPushDateGQLQuery := getMockedPullRequestLastPushDateGQLQuery(mockedPullRequest)
	mockedPullRequestLastPushDateGQLQueryBody := getMockedPullRequestLastPushDateGQLQueryBody(mockedLastCommitDate)

	tests := map[string]struct {
		inputReviewEvent string
	}{
		"when review event is APPROVE": {
			inputReviewEvent: "APPROVE",
		},
		"when review event is REQUEST_CHANGES": {
			inputReviewEvent: "REQUEST_CHANGES",
		},
		"when review event is COMMENT": {
			inputReviewEvent: "COMMENT",
		},
	}

	for _, test := range tests {
		var gotReviewEvent, gotReviewBody string

		mockedEnv := aladino.MockDefaultEnv(
			t,
			[]mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{},
				),
				mock.WithRequestMatchHandler(
					mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						rawBody, _ := ioutil.ReadAll(r.Body)

						body := ReviewRequestPostBody{}

						json.Unmarshal(rawBody, &body)

						gotReviewEvent = body.Event
						gotReviewBody = body.Body
					}),
				),
			},
			func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(aladino.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedPullRequestLastPushDateGQLQuery):
					aladino.MustWrite(w, mockedPullRequestLastPushDateGQLQueryBody)
				}
			},
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(reviewBody)}

		err := review(mockedEnv, args)

		assert.Nil(t, err)
		assert.Equal(t, test.inputReviewEvent, gotReviewEvent)
		assert.Equal(t, reviewBody, gotReviewBody)
	}
}

func TestReview_WhenBodyIsEmpty(t *testing.T) {
	mockedLastCommitDate := time.Now()

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestLastPushDateGQLQuery := getMockedPullRequestLastPushDateGQLQuery(mockedPullRequest)
	mockedPullRequestLastPushDateGQLQueryBody := getMockedPullRequestLastPushDateGQLQueryBody(mockedLastCommitDate)

	tests := map[string]struct {
		inputReviewEvent string
		wantErr          string
	}{
		"when review event is APPROVE": {
			inputReviewEvent: "APPROVE",
			wantErr:          "",
		},
		"when review event is REQUEST_CHANGES": {
			inputReviewEvent: "REQUEST_CHANGES",
			wantErr:          "review: comment required in REQUEST_CHANGES event",
		},
		"when review event is COMMENT": {
			inputReviewEvent: "COMMENT",
			wantErr:          "review: comment required in COMMENT event",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetUser,
						&github.User{
							Login: github.String("bot-account"),
						},
					),
					mock.WithRequestMatch(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						[]*github.PullRequestReview{},
					),
					mock.WithRequestMatch(
						mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
						&ReviewRequestPostBody{},
					),
				},
				func(w http.ResponseWriter, req *http.Request) {
					query := utils.MinifyQuery(aladino.MustRead(req.Body))
					switch query {
					case utils.MinifyQuery(mockedPullRequestLastPushDateGQLQuery):
						aladino.MustWrite(w, mockedPullRequestLastPushDateGQLQueryBody)
					}
				},
				aladino.MockBuiltIns(),
				nil,
			)

			args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue("")}

			gotErr := review(mockedEnv, args)

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, "review() error = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}

func TestReview_WhenPreviousAutomaticReviewWasMade(t *testing.T) {
	reviewBody := "test comment"
	reviewerLogin := "bot-account"

	mockedLastCommitDate := time.Now()
	mockedCreateDateAfterLastCommitDate := mockedLastCommitDate.Add(time.Hour)
	mockedCreateDateBeforeLastCommitDate := mockedLastCommitDate.Add(time.Hour * -1)

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockedPullRequestLastPushDateGQLQuery := getMockedPullRequestLastPushDateGQLQuery(mockedPullRequest)
	mockedPullRequestLastPushDateGQLQueryBody := getMockedPullRequestLastPushDateGQLQueryBody(mockedLastCommitDate)

	tests := map[string]struct {
		inputReviewEvent   string
		clientOptions      []mock.MockBackendOption
		shouldReviewBeDone bool
	}{
		"when last review event was APPROVE": {
			inputReviewEvent: "APPROVE",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:   github.Int64(1),
							Body: github.String(""),
							User: &github.User{
								Login: github.String(reviewerLogin),
							},
							State:       github.String("APPROVED"),
							SubmittedAt: &mockedCreateDateAfterLastCommitDate,
						},
					},
				),
			},
			shouldReviewBeDone: false,
		},
		"when last review event was REQUEST_CHANGES and the pull request has updates since then": {
			inputReviewEvent: "APPROVE",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:   github.Int64(1),
							Body: github.String(""),
							User: &github.User{
								Login: github.String(reviewerLogin),
							},
							State:       github.String("REQUEST_CHANGES"),
							SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
						},
					},
				),
			},
			shouldReviewBeDone: true,
		},
		"when last review event was REQUEST_CHANGES and the pull request has no updates since then": {
			inputReviewEvent: "APPROVE",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:   github.Int64(1),
							Body: github.String(""),
							User: &github.User{
								Login: github.String(reviewerLogin),
							},
							State:       github.String("REQUEST_CHANGES"),
							SubmittedAt: &mockedCreateDateAfterLastCommitDate,
						},
					},
				),
			},
			shouldReviewBeDone: false,
		},
		"when last review event was COMMENT and the pull request has updates since then": {
			inputReviewEvent: "APPROVE",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:   github.Int64(1),
							Body: github.String(""),
							User: &github.User{
								Login: github.String(reviewerLogin),
							},
							State:       github.String("COMMENTED"),
							SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
						},
					},
				),
			},
			shouldReviewBeDone: true,
		},
		"when last review event was COMMENT and the pull request has no updates since then": {
			inputReviewEvent: "APPROVE",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetUser,
					&github.User{
						Login: github.String("bot-account"),
					},
				),
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:   github.Int64(1),
							Body: github.String(""),
							User: &github.User{
								Login: github.String(reviewerLogin),
							},
							State:       github.String("COMMENTED"),
							SubmittedAt: &mockedCreateDateAfterLastCommitDate,
						},
					},
				),
			},
			shouldReviewBeDone: false,
		},
	}

	for _, test := range tests {
		var isPostReviewRequestPerformed bool
		mockedEnv := aladino.MockDefaultEnv(
			t,
			append(
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetUser,
						&github.User{
							Login: github.String("bot-account"),
						},
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(mockedPullRequest))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.PostReposPullsReviewsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							isPostReviewRequestPerformed = true
						}),
					),
				},
				test.clientOptions...,
			),
			func(w http.ResponseWriter, req *http.Request) {
				query := utils.MinifyQuery(aladino.MustRead(req.Body))
				switch query {
				case utils.MinifyQuery(mockedPullRequestLastPushDateGQLQuery):
					aladino.MustWrite(w, mockedPullRequestLastPushDateGQLQueryBody)
				}
			},
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{aladino.BuildStringValue(test.inputReviewEvent), aladino.BuildStringValue(reviewBody)}

		err := review(mockedEnv, args)

		assert.Nil(t, err)
		assert.Equal(t, test.shouldReviewBeDone, isPostReviewRequestPerformed)
	}
}

func getMockedPullRequestLastPushDateGQLQuery(mockedPullRequest *github.PullRequest) string {
	mockedPullRequestNumber := host.GetPullRequestNumber(mockedPullRequest)
	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)

	return fmt.Sprintf(`{
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
	}`, mockedPullRequestNumber, mockRepo, mockOwner)
}

func getMockedPullRequestLastPushDateGQLQueryBody(mockedLastCommitDate time.Time) string {
	return fmt.Sprintf(`{
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
}
