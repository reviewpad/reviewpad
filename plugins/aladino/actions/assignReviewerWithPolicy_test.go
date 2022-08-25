// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignReviewerWithPolicy = plugins_aladino.PluginBuiltIns().Actions["assignReviewerWithPolicy"].Code

func TestAssignReviewerWithPolicy_WhenTotalRequiredReviewersIsZero(t *testing.T) {
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
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewerWithPolicy: total required reviewers can't be 0")
}

func TestAssignReviewerWithPolicy_WhenListOfReviewersIsEmpty(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{}), aladino.BuildIntValue(1), aladino.BuildStringValue("random")}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.EqualError(t, err, "assignReviewerWithPolicy: list of reviewers can't be empty")
}

func TestAssignReviewerWithPolicy_WhenAuthorIsInListOfReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(authorLogin),
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "pr author shouldn't be assigned as a reviewer")
}

func TestAssignReviewerWithPolicy_WhenTotalRequiredReviewersIsMoreThanTotalAvailableReviewers(t *testing.T) {
	var gotReviewers []string
	reviewerLogin := "mary"
	authorLogin := "john"
	totalRequiredReviewers := 2
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(totalRequiredReviewers),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "the list of assign reviewers should be all provided reviewers")
}

func TestAssignReviewerWithPolicy_WhenListReviewsRequestFails(t *testing.T) {
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
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestAssignReviewerWithPolicy_WhenPullRequestAlreadyHasReviews(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{
		reviewerLogin,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
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
						State: github.String("COMMENTED"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewerWithPolicy_WhenPullRequestAlreadyHasApproval(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerLogin := "mary"
	wantReviewers := []string{}

	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
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
						State: github.String("APPROVED"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerLogin),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewerWithPolicy_WhenPullRequestAlreadyHasRequestedReviewers(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	wantReviewers := []string{
		reviewerB,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerA)},
		},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerA),
				aladino.BuildStringValue(reviewerB),
			},
		),
		aladino.BuildIntValue(2),
		aladino.BuildStringValue("random"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers, "when a reviewer already has a requested review, then it shouldn't be re-requested")
}

// Test scenario description:
// The mocked pull request has an assigned reviewer ("mary") which hasn't made any review yet
// The provided reviewers list contains an already assigned reviewer ("mary")
// Since a review has already been requested to the already assigned reviewer, so there's no available reviewers left
func TestAssignReviewerWithPolicy_HasNoAvailableReviewers(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerLogin := "mary"
	totalRequiredReviewers := 1
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerLogin)},
		},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
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
		nil,
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
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}

func TestAssignReviewerWithPolicy_WhenPullRequestAlreadyApproved(t *testing.T) {
	var isRequestReviewersRequestPerformed bool
	authorLogin := "john"
	reviewerA := "mary"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{
			{Login: github.String(reviewerA)},
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
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
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
		nil,
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
	err := assignReviewerWithPolicy(mockedEnv, args)
	assert.Nil(t, err)
	assert.False(t, isRequestReviewersRequestPerformed, "the action shouldn't request for reviewers")
}

func TestAssignReviewerWithPolicy_WhenRoundRobinPolicy(t *testing.T) {
	var gotReviewers []string
	prNum := 258
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	wantReviewers := []string{
		reviewerA,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Number:             &prNum,
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{
					{
						ID:   github.Int64(1),
						Body: github.String(""),
						User: &github.User{
							Login: github.String(reviewerA),
						},
						State: github.String("COMMENTED"),
					},
					{
						ID:   github.Int64(2),
						Body: github.String(""),
						User: &github.User{
							Login: github.String(reviewerB),
						},
						State: github.String("COMMENTED"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerA),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("round-robin"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}

func TestAssignReviewerWithPolicy_WhenReviewpadPolicy(t *testing.T) {
	var gotReviewers []string
	authorLogin := "john"
	reviewerA := "mary"
	reviewerB := "steve"
	wantReviewers := []string{
		reviewerB,
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User:               &github.User{Login: github.String(authorLogin)},
		RequestedReviewers: []*github.User{},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatch(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				[]*github.PullRequestReview{},
			),
			mock.WithRequestMatchHandler(
				mock.GetSearchIssues,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					var issuesCount int
					if strings.Contains(r.URL.Query()["q"][0], reviewerA) {
						issuesCount = 20
					} else {
						issuesCount = 12
					}
					w.Write(mock.MustMarshal(github.IssuesSearchResult{
						Total: &issuesCount,
					}))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.PostReposPullsRequestedReviewersByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := ReviewersRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotReviewers = body.Reviewers
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
				aladino.BuildStringValue(reviewerA),
				aladino.BuildStringValue(reviewerB),
			},
		),
		aladino.BuildIntValue(1),
		aladino.BuildStringValue("reviewpad"),
	}
	err := assignReviewerWithPolicy(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantReviewers, gotReviewers)
}
