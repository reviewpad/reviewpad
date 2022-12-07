// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var size = plugins_aladino.PluginBuiltIns().Functions["size"].Code

func TestSize_WhenRegexMatchFails(t *testing.T) {
	tests := map[string]struct {
		args    []aladino.Value
		wantVal aladino.Value
	}{
		"no exclusions": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})},
			wantVal: aladino.BuildIntValue(6),
		},
		"exclude yml": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("*.yml")})},
			wantVal: aladino.BuildIntValue(3),
		},
		"invalid regex": {
			args:    []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("*[.yml")})},
			wantVal: aladino.BuildIntValue(6),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			changes := 3
			mockedPullRequestFileList := &[]*github.CommitFile{
				{
					Filename: github.String("reviewpad.yml"),
					Patch:    nil,
					Changes:  &changes,
				},
				{
					Filename: github.String("main.go"),
					Patch:    nil,
					Changes:  &changes,
				},
			}
			additions := 6
			deletions := 0
			mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
				Additions: github.Int(additions),
				Deletions: github.Int(deletions),
			})

			mockedEnv := aladino.MockDefaultEnv(
				t,
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
						}),
					),
					mock.WithRequestMatchHandler(
						mock.GetReposPullsByOwnerByRepoByPullNumber,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
						}),
					),
				},
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			gotSize, err := size(mockedEnv, test.args)

			assert.Nil(t, err)
			assert.Equal(t, test.wantVal, gotSize)
		})
	}
}
