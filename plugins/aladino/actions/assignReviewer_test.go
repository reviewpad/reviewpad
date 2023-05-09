// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var assignReviewer = plugins_aladino.PluginBuiltIns().Actions["assignReviewer"].Code

type reviewers struct {
	Reviewers []string `json:"reviewers"`
}

func TestAssignReviewer_WhenPolicyIsNotValid(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	invalidPolicy := "INVALID_POLICY"
	allowedPolicies := map[string]bool{"random": true, "round-robin": true, "reviewpad": true}

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
			},
		),
		aladino.BuildIntValue(0),
		aladino.BuildStringValue(invalidPolicy),
	}
	err := assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, fmt.Sprintf("assignReviewer: policy %s is not supported. allowed policies %v", invalidPolicy, allowedPolicies))
}

func TestAssignReviewer_WhenTotalRequiredReviewersIsZero(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
			},
		),
		aladino.BuildIntValue(0),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: total required reviewers can't be 0")
}

func TestAssignReviewer_WhenListOfReviewersIsEmpty(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{
		aladino.BuildArrayValue([]aladino.Value{}),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewer: list of reviewers can't be empty")
}

func TestAssignReviewer_WhenAuthorIsInListOfReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author:             &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-11-11T13:36:05Z",
							  "committedDate": "2022-11-11T13:36:05Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(authorLogin),
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "pr author shouldn't be assigned as a reviewer")
}

func TestAssignReviewer_WhenTotalRequiredReviewersIsMoreThanTotalAvailableReviewers(t *testing.T) {
	var gotReviewers []string
	reviewerLogin := "mary"
	authorLogin := "john"
	totalRequiredReviewers := 2
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author:             &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{},
	})
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-11-11T13:36:05Z",
							  "committedDate": "2022-11-11T13:36:05Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "the list of assign reviewers should be all provided reviewers")
}

func TestAssignReviewer_WhenListReviewsRequestFails(t *testing.T) {
	failMessage := "ListReviewsRequestFail"
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
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue("jane"),
			},
		),
		aladino.BuildIntValue(3),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignReviewer_WhenPullRequestAlreadyHasReviews(t *testing.T) {
	time := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author:             &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
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
						SubmittedAt: &github.Timestamp{Time: time},
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewer_WhenPullRequestAlreadyHasApproval(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{}
	time := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)

	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author:             &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
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
						SubmittedAt: &github.Timestamp{Time: time},
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewer_WhenPullRequestAlreadyHasRequestedReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	wantReviewers := []string{
		reviewerB,
	}
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{
			Users: []*pbc.User{
				{
					Login: reviewerA,
				},
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerA),
				aladino.BuildStringValue(reviewerB),
			},
		),
		aladino.BuildIntValue(2),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "when a reviewer already has a requested review, then it shouldn't be re-requested")
}

// Test scenario description:
// The mocked pull request has an assigned reviewer ("mary") which hasn't made any review yet
// The provided reviewers list contains an already assigned reviewer ("mary")
// Since a review has already been requested to the already assigned reviewer, so there's no available reviewers left
func TestAssignReviewer_HasNoAvailableReviewers(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerLogin := "mary"
	totalRequiredReviewers := 1
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{
			Login: authorLogin,
		},
		RequestedReviewers: &pbc.RequestedReviewers{
			Users: []*pbc.User{
				{
					Login: reviewerLogin,
				},
			},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the request reviewers request was performed then the reviewers were assigned to the pull request
					isRequestReviewersRequestPerformed = true
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}

func TestAssignReviewer_WhenPullRequestAlreadyApproved(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerA := "mary"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{
			Users: []*pbc.User{
				{Login: reviewerA},
			},
		},
	})
	reviewID := int64(1)
	reviewState := "APPROVED"
	time := time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	mockedReviews := []*github.PullRequestReview{
		{
			ID:          &reviewID,
			State:       &reviewState,
			Body:        github.String(""),
			User:        &github.User{Login: github.String(reviewerA)},
			SubmittedAt: &github.Timestamp{Time: time},
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				mockedReviews,
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the request reviewers request was performed then the reviewers were assigned to the pull request
					isRequestReviewersRequestPerformed = true
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerA),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)
	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}

func TestAssignReviewer_ReRequest(t *testing.T) {
	mockedReviewerLogin := "mary"
	tests := map[string]struct {
		shouldRequestReview  bool
		mockedReviews        func() []*github.PullRequestReview
		lastPushDateResponse string
	}{
		"when reviewer last review is approved": {
			shouldRequestReview: false,
			mockedReviews: func() []*github.PullRequestReview {
				time := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("APPROVED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: time},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer last review is commented": {
			shouldRequestReview: true,
			mockedReviews: func() []*github.PullRequestReview {
				time := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("COMMENTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: time},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer last review is changes requested": {
			shouldRequestReview: true,
			mockedReviews: func() []*github.PullRequestReview {
				time := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("CHANGES_REQUESTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: time},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer has approved and then requested changes": {
			shouldRequestReview: true,
			mockedReviews: func() []*github.PullRequestReview {
				timeApproved := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				timeRequestedChanges := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("APPROVED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeApproved},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
					{
						ID:          github.Int64(2),
						State:       github.String("CHANGES_REQUESTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeRequestedChanges},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer has requested changes and then approved": {
			shouldRequestReview: false,
			mockedReviews: func() []*github.PullRequestReview {
				timeApproved := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
				timeRequestedChanges := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("CHANGES_REQUESTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeRequestedChanges},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
					{
						ID:          github.Int64(2),
						State:       github.String("APPROVED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeApproved},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer has requested changes and no commits have been pushed after": {
			shouldRequestReview: false,
			mockedReviews: func() []*github.PullRequestReview {
				timeRequestedChanges := time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("CHANGES_REQUESTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeRequestedChanges},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-01-11T01:01:01Z",
								  "committedDate": "2022-01-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
		"when reviewer has requested changes and some commits have been pushed after": {
			shouldRequestReview: true,
			mockedReviews: func() []*github.PullRequestReview {
				timeRequestedChanges := time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC)
				return []*github.PullRequestReview{
					{
						ID:          github.Int64(1),
						State:       github.String("CHANGES_REQUESTED"),
						Body:        github.String(""),
						SubmittedAt: &github.Timestamp{Time: timeRequestedChanges},
						User:        &github.User{Login: github.String(mockedReviewerLogin)},
					},
				}
			},
			lastPushDateResponse: `
			{
				"data": {
					"repository": {
						"pullRequest": {
						  "timelineItems": {
							"nodes": [
							  {
								"__typename": "PullRequestCommit",
								"commit": {
								  "pushedDate": "2022-03-11T01:01:01Z",
								  "committedDate": "2022-03-11T01:01:01Z"
								}
							  }
							]
						  }
						}
					}
				}
			}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			isRequestReviewersRequestPerformed := false
			mockedAuthorLogin := "john"
			mockedReviews := test.mockedReviews()
			mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
				Author: &pbc.User{Login: mockedAuthorLogin},
				RequestedReviewers: &pbc.RequestedReviewers{
					Users: []*pbc.User{
						{Login: mockedReviewerLogin},
					},
				},
			})
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatch(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						mockedReviews,
					),
					mock.WithRequestMatchHandler(
						mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							// If the request reviewers request was performed then the reviewers were assigned to the pull request
							isRequestReviewersRequestPerformed = true
						}),
					),
				},
				func(w http.ResponseWriter, r *http.Request) {
					utils.MustWrite(w, test.lastPushDateResponse)
				},
				mockedCodeReview,
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			)
			args := []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue(mockedReviewerLogin),
					},
				),
				aladino.BuildIntValue(1),
				aladino.BuildStringValue("random"),
			}
			err := assignReviewer(mockedEnv, args)

			assert.Nil(t, err)
			assert.Equal(t, test.shouldRequestReview, isRequestReviewersRequestPerformed)
		})
	}
}

func TestAssignReviewer_WithPolicy(t *testing.T) {
	mockedReviewerALogin := "mary"
	mockedReviewerBLogin := "peter"
	mockedReviewerCLogin := "jeff"
	tests := map[string]struct {
		inputPolicy       string
		clientOptions     []mock.MockBackendOption
		expectedReviewers []string
		wantErr           string
	}{
		"when policy is round-robin": {
			inputPolicy:       "round-robin",
			clientOptions:     []mock.MockBackendOption{},
			expectedReviewers: []string{mockedReviewerCLogin},
			wantErr:           "",
		},
		"when policy is reviewpad": {
			inputPolicy: "reviewpad",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchPages(
					mock.GetSearchIssues,
					&github.IssuesSearchResult{
						Total: github.Int(0),
					},
					&github.IssuesSearchResult{
						Total: github.Int(0),
					},
					&github.IssuesSearchResult{
						Total: github.Int(0),
					},
				),
			},
			expectedReviewers: []string{mockedReviewerALogin},
			wantErr:           "",
		},
		"when policy is reviewpad and request for search issues fails": {
			inputPolicy: "reviewpad",
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetSearchIssues,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mock.WriteError(
							w,
							http.StatusInternalServerError,
							"SearchIssuesRequestFail",
						)
					}),
				),
			},
			expectedReviewers: []string{},
			wantErr:           "SearchIssuesRequestFail",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requestedReviewers := []string{}
			mockedAuthorLogin := "john"
			mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
				Author:             &pbc.User{Login: mockedAuthorLogin},
				RequestedReviewers: &pbc.RequestedReviewers{},
			})
			mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
				t,
				append(
					[]mock.MockBackendOption{
						mock.WithRequestMatch(
							mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
							[]*github.PullRequestReview{},
						),
						mock.WithRequestMatchHandler(
							mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
							http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								rawBody, _ := io.ReadAll(r.Body)
								body := reviewers{}

								utils.MustUnmarshal(rawBody, &body)

								requestedReviewers = body.Reviewers
							}),
						),
					},
					test.clientOptions...,
				),
				func(w http.ResponseWriter, r *http.Request) {
					utils.MustWrite(w, `
					{
						"data": {
						  "repository": {
							"pullRequest": {
							  "timelineItems": {
								"nodes": [
								  {
									"__typename": "PullRequestCommit",
									"commit": {
									  "pushedDate": "2022-01-11T01:01:01Z",
									  "committedDate": "2022-01-11T01:01:01Z"
									}
								  }
								]
							  }
							}
						  }
						}
					  }
					`)
				},
				mockedCodeReview,
				aladino.GetDefaultPullRequestFileList(),
				aladino.MockBuiltIns(),
				nil,
			)
			args := []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue(mockedReviewerALogin),
						aladino.BuildStringValue(mockedReviewerBLogin),
						aladino.BuildStringValue(mockedReviewerCLogin),
					},
				),
				aladino.BuildIntValue(1),
				aladino.BuildStringValue(test.inputPolicy),
			}

			gotErr := assignReviewer(mockedEnv, args)

			if gotErr != nil && gotErr.(*github.ErrorResponse).Message != test.wantErr {
				assert.FailNow(t, "assignReviewer() error = %v, wantErr %v", gotErr, test.wantErr)
			}
			assert.Equal(t, test.expectedReviewers, requestedReviewers)
		})
	}
}

func TestAssignReviewer_WhenPullRequestAlreadyApprovedBy2Reviewers(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: authorLogin},
		RequestedReviewers: &pbc.RequestedReviewers{
			Users: []*pbc.User{
				{Login: reviewerA},
			},
		},
	})
	reviewID := int64(1)
	reviewState := "APPROVED"
	mockedReviews := []*github.PullRequestReview{
		{
			ID:    &reviewID,
			State: &reviewState,
			Body:  github.String(""),
			User:  &github.User{Login: github.String(reviewerA)},
		},
		{
			ID:    &reviewID,
			State: &reviewState,
			Body:  github.String(""),
			User:  &github.User{Login: github.String(reviewerB)},
		},
	}
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				mockedReviews,
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// If the request reviewers request was performed then the reviewers were assigned to the pull request
					isRequestReviewersRequestPerformed = true
				}),
			),
		},
		func(w http.ResponseWriter, r *http.Request) {
			utils.MustWrite(w, `
			{
				"data": {
				  "repository": {
					"pullRequest": {
					  "timelineItems": {
						"nodes": [
						  {
							"__typename": "PullRequestCommit",
							"commit": {
							  "pushedDate": "2022-01-11T01:01:01Z",
							  "committedDate": "2022-01-11T01:01:01Z"
							}
						  }
						]
					  }
					}
				  }
				}
			  }
			`)
		},
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	args := []aladino.Value{
		aladino.BuildArrayValue(
			[]aladino.Value{
				aladino.BuildStringValue(reviewerA),
				aladino.BuildStringValue(reviewerB),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewer(mockedEnv, args)
	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}
