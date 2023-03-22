// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var deleteHeadBranch = plugins_aladino.PluginBuiltIns().Actions["deleteHeadBranch"].Code

func TestDeleteHeadBranch(t *testing.T) {
	now := time.Now()
	isDeleteHeadBranchRequestPerformed := false
	tests := map[string]struct {
		clientOptions           []mock.MockBackendOption
		codeReview              *pbe.CodeReview
		graphQLHandler          http.HandlerFunc
		deleteShouldBePerformed bool
		err                     error
	}{
		"when pull request is closed": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: false,
				ClosedAt: timestamppb.New(now),
			}),
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						isDeleteHeadBranchRequestPerformed = true
						w.WriteHeader(http.StatusNoContent)
					}),
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"ref": {
								"id": "XYSD9fbcCu"
							}
						}
					}
				}`)
			},
			deleteShouldBePerformed: true,
		},
		"when pull request is merged": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: true,
				ClosedAt: timestamppb.New(now),
			}),
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
						isDeleteHeadBranchRequestPerformed = true
					}),
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"ref": {
								"id": "XYSD9fbcCu"
							}
						}
					}
				}`)
			},
			deleteShouldBePerformed: true,
		},
		"when pull request is not merged or closed": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: false,
				ClosedAt: nil,
			}),
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusNoContent)
					}),
				),
			},
			err: nil,
		},
		"when delete head branch ref request fails": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: true,
			}),
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.DeleteReposGitRefsByOwnerByRepoByRef,
					http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						utils.MustWrite(w, `{
							"message": "Reference does not exist",
							"documentation_url": "https://docs.github.com/rest/reference/git#delete-a-reference"
						}`)
					}),
				),
			},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"ref": {
								"id": "XYSD9fbcCu"
							}
						}
					}
				}`)
			},
			err: &github.ErrorResponse{
				Message:          "Reference does not exist",
				DocumentationURL: "https://docs.github.com/rest/reference/git#delete-a-reference",
			},
		},
		"when pull request is from fork": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: true,
				Head: &pbe.Branch{
					Repo: &pbe.Repository{
						IsFork: true,
					},
				},
			}),
			clientOptions: []mock.MockBackendOption{},
			err:           nil,
		},
		"when head branch doesn't exist": {
			codeReview: aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
				IsMerged: false,
			}),
			clientOptions: []mock.MockBackendOption{},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
						"repository": {
							"ref": null
						}
					}
				}`)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnvWithCodeReview(
				t,
				test.clientOptions,
				test.graphQLHandler,
				test.codeReview,
				aladino.MockBuiltIns(),
				nil,
			)

			err := deleteHeadBranch(mockedEnv, []aladino.Value{})

			// this allows simplified checking of github error response equality
			if e, ok := err.(*github.ErrorResponse); ok {
				e.Response = nil
			}

			assert.Equal(t, test.err, err)

			assert.Equal(t, test.deleteShouldBePerformed, isDeleteHeadBranchRequestPerformed)

			isDeleteHeadBranchRequestPerformed = false
		})
	}
}
