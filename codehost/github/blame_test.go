// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetGitBlame(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		graphQLHandler http.HandlerFunc
		owner          string
		repo           string
		commitSHA      string
		filePaths      []string
		wantRes        *github.GitBlame
		wantErr        error
	}{
		"when graphql errors": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			},
			wantErr: fmt.Errorf("error executing blame query: Message: 403 Forbidden; body: \"\", Locations: []"),
		},
		"when blame information is not found": {
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{"data": {}}`)
			},
			wantErr: errors.New("error getting blame information: no blame information found"),
		},
		"when successful": {
			commitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
			filePaths: []string{"reviewpad.yml"},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
					  "repository": {
						"object": {
						  "blame": {
							"ranges": [
							  {
								"startingLine": 1,
								"endingLine": 10,
								"age": 10,
								"commit": {
								  "author": {
									"user": {
									  "login": "reviewpad"
									}
								  }
								}
							  }
							]
						  }
						}
					  }
					}
				}`)
			},
			wantRes: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 10,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "reviewpad",
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := aladino.MockDefaultGithubClient(nil, test.graphQLHandler)
			res, err := client.GetGitBlame(ctx, test.owner, test.repo, test.commitSHA, test.filePaths)
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantRes, res)
		})
	}
}
