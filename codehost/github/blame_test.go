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
						  "blame0": {
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
		"when successful with empty ranges": {
			commitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
			filePaths: []string{"reviewpad.yml"},
			graphQLHandler: func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{
					"data": {
					  "repository": {
						"object": {
						  "blame0": {
							"ranges": [
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
						LineCount: 0,
						Lines:     []github.GitBlameFileLines{},
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

func TestSliceGitBlame(t *testing.T) {
	tests := map[string]struct {
		blame    *github.GitBlame
		filePath string
		slices   []github.FileSliceLimits
		wantRes  *github.GitBlame
		wantErr  error
	}{
		"when file does not exist in blame": {
			blame: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 100,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   100,
								Author:   "reviewpad[bot]",
							},
						},
					},
				},
			},
			wantErr:  errors.New("file not found in blame information"),
			filePath: "README.md",
		},
		"when line is found": {
			blame: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 100,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "reviewpad[bot]",
							},
							{
								FromLine: 11,
								ToLine:   30,
								Author:   "test",
							},
							{
								FromLine: 31,
								ToLine:   100,
								Author:   "reviewpad[bot]",
							},
						},
					},
				},
			},
			wantRes: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 100,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 20,
								ToLine:   25,
								Author:   "test",
							},
						},
					},
				},
			},
			slices: []github.FileSliceLimits{
				{
					From: 20,
					To:   25,
				},
			},
			filePath: "reviewpad.yml",
		},
		"when line is out of range": {
			blame: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 100,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "reviewpad[bot]",
							},
							{
								FromLine: 11,
								ToLine:   30,
								Author:   "test",
							},
							{
								FromLine: 31,
								ToLine:   100,
								Author:   "reviewpad[bot]",
							},
						},
					},
				},
			},
			wantRes: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"reviewpad.yml": {
						FilePath:  "reviewpad.yml",
						LineCount: 100,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 10,
								ToLine:   10,
								Author:   "reviewpad[bot]",
							},
							{
								FromLine: 11,
								ToLine:   30,
								Author:   "test",
							},
							{
								FromLine: 31,
								ToLine:   45,
								Author:   "reviewpad[bot]",
							},
						},
					},
				},
			},
			slices: []github.FileSliceLimits{
				{
					From: 10,
					To:   45,
				},
			},
			filePath: "reviewpad.yml",
		},
		"when multiple files in blame": {
			blame: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"bar.go": {
						FilePath:  "bar.go",
						LineCount: 20,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "peter",
							},
							{
								FromLine: 11,
								ToLine:   20,
								Author:   "maria",
							},
						},
					},
					"foo.go": {
						FilePath:  "foo.go",
						LineCount: 10,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "john",
							},
						},
					},
				},
			},
			wantRes: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"bar.go": {
						FilePath:  "bar.go",
						LineCount: 20,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 10,
								ToLine:   10,
								Author:   "peter",
							},
							{
								FromLine: 11,
								ToLine:   15,
								Author:   "maria",
							},
						},
					},
				},
			},
			slices: []github.FileSliceLimits{
				{
					From: 10,
					To:   15,
				},
			},
			filePath: "bar.go",
		},
		"when multiple slices": {
			blame: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"bar.go": {
						FilePath:  "bar.go",
						LineCount: 50,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "peter",
							},
							{
								FromLine: 11,
								ToLine:   20,
								Author:   "maria",
							},
							{
								FromLine: 21,
								ToLine:   30,
								Author:   "jack",
							},
							{
								FromLine: 31,
								ToLine:   35,
								Author:   "jill",
							},
							{
								FromLine: 36,
								ToLine:   45,
								Author:   "reviewpad[bot]",
							},
							{
								FromLine: 46,
								ToLine:   50,
								Author:   "james",
							},
						},
					},
					"foo.go": {
						FilePath:  "foo.go",
						LineCount: 10,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 1,
								ToLine:   10,
								Author:   "john",
							},
						},
					},
				},
			},
			wantRes: &github.GitBlame{
				CommitSHA: "41a7b4350148def9e98c9352a64cc3ac95c1de9c",
				Files: map[string]github.GitBlameFile{
					"bar.go": {
						FilePath:  "bar.go",
						LineCount: 50,
						Lines: []github.GitBlameFileLines{
							{
								FromLine: 10,
								ToLine:   10,
								Author:   "peter",
							},
							{
								FromLine: 11,
								ToLine:   15,
								Author:   "maria",
							},
							{
								FromLine: 31,
								ToLine:   35,
								Author:   "jill",
							},
							{
								FromLine: 36,
								ToLine:   45,
								Author:   "reviewpad[bot]",
							},
						},
					},
				},
			},
			slices: []github.FileSliceLimits{
				{
					From: 10,
					To:   15,
				},
				{
					From: 31,
					To:   45,
				},
			},
			filePath: "bar.go",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := aladino.MockDefaultGithubClient(nil, nil)
			res, err := client.SliceGitBlame(test.blame, test.filePath, test.slices)
			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantRes, res)
		})
	}
}
