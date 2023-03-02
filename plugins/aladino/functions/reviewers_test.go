// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var reviewers = plugins_aladino.PluginBuiltIns().Functions["reviewers"].Code

func TestReviewers_WhenListReviewsRequestFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						"ListReviewsRequestFail",
					)
				}),
			),
		},
		nil, aladino.MockBuiltIns(),
		nil,
	)

	gotReviewers, gotErr := reviewers(mockedEnv, []aladino.Value{})

	assert.Equal(t, "ListReviewsRequestFail", gotErr.(*github.ErrorResponse).Message)
	assert.Nil(t, gotReviewers)
}

func TestReviewers(t *testing.T) {
	tests := map[string]struct {
		clientOptions []mock.MockBackendOption
		wantReviewers aladino.Value
		wantErr       string
	}{
		"when pull request has no reviews": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{},
				),
			},
			wantReviewers: aladino.BuildArrayValue([]aladino.Value{}),
			wantErr:       "",
		},
		"when pull request has reviews": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:    github.Int64(1),
							Body:  github.String("Here is the body for the review."),
							State: github.String("APPROVED"),
							User: &github.User{
								Login: github.String("mary"),
							},
						},
					},
				),
			},
			wantReviewers: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("mary")}),
			wantErr:       "",
		},
		"when pull request has more than one review of the same user": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:    github.Int64(1),
							Body:  github.String("Here is the body for the review."),
							State: github.String("REQUESTED_CHANGES"),
							User: &github.User{
								Login: github.String("mary"),
							},
						},
						{
							ID:    github.Int64(2),
							Body:  github.String("Here is the body for the review."),
							State: github.String("APPROVED"),
							User: &github.User{
								Login: github.String("mary"),
							},
						},
					},
				),
			},
			wantReviewers: aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("mary")}),
			wantErr:       "",
		},
		"when pull request has more than one review of different users": {
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsReviewsByOwnerByRepoByPullNumber,
					[]*github.PullRequestReview{
						{
							ID:    github.Int64(1),
							Body:  github.String("Here is the body for the review."),
							State: github.String("REQUESTED_CHANGES"),
							User: &github.User{
								Login: github.String("mary"),
							},
						},
						{
							ID:    github.Int64(2),
							Body:  github.String("Here is the body for the review."),
							State: github.String("APPROVED"),
							User: &github.User{
								Login: github.String("john"),
							},
						},
					},
				),
			},
			wantReviewers: aladino.BuildArrayValue(
				[]aladino.Value{
					aladino.BuildStringValue("mary"),
					aladino.BuildStringValue("john"),
				},
			),
			wantErr: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(t, test.clientOptions, nil, aladino.MockBuiltIns(), nil)

			gotReviewers, gotErr := reviewers(mockedEnv, []aladino.Value{})

			if gotErr != nil && gotErr.Error() != test.wantErr {
				assert.FailNow(t, fmt.Sprintf("reviewers() error = %v, wantErr %v", gotErr, test.wantErr))
			}
			assert.Equal(t, test.wantReviewers, gotReviewers)
		})
	}
}
