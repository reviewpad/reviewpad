// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	host "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var isWaitingForReview = plugins_aladino.PluginBuiltIns().Functions["isWaitingForReview"].Code

func TestIsWaitingForReview_WhenRequestFail(t *testing.T) {
	mockedGithubClient := aladino.MockDefaultGithubClient(
		nil,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		},
	)
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetails()

	mockOwner := host.GetPullRequestBaseOwnerName(mockedPullRequest)
	mockRepo := host.GetPullRequestBaseRepoName(mockedPullRequest)
	mockPrNumber := host.GetPullRequestNumber(mockedPullRequest)

	date, err := mockedGithubClient.GetPullRequestLastPushDate(context.Background(), mockOwner, mockRepo, mockPrNumber)

	assert.NotNil(t, err)
	assert.Equal(t, time.Time{}, date)
}

func TestIsWaitingForReview_WhenNoCommits(t *testing.T) {
	mockedCommits := []*github.RepositoryCommit{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(
						aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
							RequestedReviewers: []*github.User{},
							RequestedTeams:     []*github.Team{},
						}),
					))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedCommits))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)
	args := []aladino.Value{}
	gotValue, err := isWaitingForReview(mockedEnv, args)

	assert.Nil(t, err, "Expected no error")
	assert.True(t, gotValue.Equals(aladino.BuildBoolValue(false)))
}

func TestIsWaitingForReview_WhenHasNoReviews(t *testing.T) {
	mockedLastCommitDate := time.Now()
	tests := map[string]struct {
		requestedReviewers []*github.User
		requestedTeams     []*github.Team
		wantValue          aladino.Value
	}{
		"from user": {
			requestedReviewers: []*github.User{{
				Login: github.String("john"),
			}},
			requestedTeams: []*github.Team{},
			wantValue:      aladino.BuildBoolValue(true),
		},
		"from team": {
			requestedReviewers: []*github.User{},
			requestedTeams: []*github.Team{{
				Name: github.String("team"),
			}},
			wantValue: aladino.BuildBoolValue(true),
		},
		"from user and team": {
			requestedReviewers: []*github.User{{
				Login: github.String("john"),
			}},
			requestedTeams: []*github.Team{{
				Name: github.String("team"),
			}},
			wantValue: aladino.BuildBoolValue(true),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(
								aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
									RequestedReviewers: test.requestedReviewers,
									RequestedTeams:     test.requestedTeams,
								}),
							))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write(mock.MustMarshal(
								[]*github.RepositoryCommit{{
									Commit: &github.Commit{
										Committer: &github.CommitAuthor{
											Date: &mockedLastCommitDate,
										},
									},
								}},
							))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write(mock.MustMarshal([]*github.PullRequestReview{}))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)
			args := []aladino.Value{}
			gotValue, err := isWaitingForReview(mockedEnv, args)

			assert.Nil(t, err, "Expected no error")
			assert.True(t, gotValue.Equals(test.wantValue), "Expected %v, got %v", test.wantValue, gotValue)
		})
	}
}

func TestIsWaitingForReview_WhenHasReviews(t *testing.T) {
	t.Skip("FIXME: #331")
	mockedAuthorLogin := "author"
	mockedLastCommitDate := time.Now()
	mockedCreateDateAfterLastCommitDate := mockedLastCommitDate.Add(time.Hour)
	mockedCreateDateBeforeLastCommitDate := mockedLastCommitDate.Add(time.Hour * -1)
	mockedCommits := []*github.RepositoryCommit{{
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: &mockedLastCommitDate,
			},
		},
	}}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		RequestedReviewers: []*github.User{},
		RequestedTeams:     []*github.Team{},
		User:               &github.User{Login: github.String(mockedAuthorLogin)},
	})

	tests := map[string]struct {
		reviewers []*github.PullRequestReview
		wantValue aladino.Value
	}{
		"from nil user": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: nil,
				},
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from author": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String(mockedAuthorLogin),
				},
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from user but outdated CHANGES_REQUESTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("CHANGES_REQUESTED"),
				SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(true),
		},
		"from user but outdated COMMENTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("COMMENTED"),
				SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(true),
		},
		"from user but outdated APPROVED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("APPROVED"),
				SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from user and up to date CHANGES_REQUESTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("CHANGES_REQUESTED"),
				SubmittedAt: &mockedCreateDateAfterLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from user and up to date COMMENTED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("COMMENTED"),
				SubmittedAt: &mockedCreateDateAfterLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from user and up to date APPROVED": {
			reviewers: []*github.PullRequestReview{{
				User: &github.User{
					Login: github.String("user"),
				},
				State:       github.String("APPROVED"),
				SubmittedAt: &mockedCreateDateAfterLastCommitDate,
			}},
			wantValue: aladino.BuildBoolValue(false),
		},
		"from user and more than one": {
			reviewers: []*github.PullRequestReview{
				{
					User: &github.User{
						Login: github.String("user"),
					},
					State:       github.String("CHANGES_REQUESTED"),
					SubmittedAt: &mockedCreateDateBeforeLastCommitDate,
				},
				{
					User: &github.User{
						Login: github.String("user"),
					},
					State:       github.String("APPROVED"),
					SubmittedAt: &mockedCreateDateAfterLastCommitDate,
				},
			},
			wantValue: aladino.BuildBoolValue(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							w.Write(mock.MustMarshal(mockedPullRequest))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write(mock.MustMarshal(mockedCommits))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							w.Write(mock.MustMarshal(test.reviewers))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)
			args := []aladino.Value{}
			gotValue, err := isWaitingForReview(mockedEnv, args)

			assert.Nil(t, err, "Expected no error")
			assert.True(t, gotValue.Equals(test.wantValue), "Expected %v, got %v", test.wantValue, gotValue)
		})
	}
}
