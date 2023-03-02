// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var isMerged = plugins_aladino.PluginBuiltIns().Functions["isMerged"].Code

func TestIsMerged(t *testing.T) {
	tests := map[string]struct {
		backendOptions []mock.MockBackendOption
		wantResult     aladino.Value
		wantErr        error
	}{
		"when pull request is merged": {
			backendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
						Merged: github.Bool(true),
					}),
				),
			},
			wantResult: aladino.BuildBoolValue(true),
		},
		"when pull request is not merged": {
			backendOptions: []mock.MockBackendOption{
				mock.WithRequestMatch(
					mock.GetReposPullsByOwnerByRepoByPullNumber,
					aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
						Merged: github.Bool(false),
					}),
				),
			},
			wantResult: aladino.BuildBoolValue(false),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockedEnv := aladino.MockDefaultEnv(
				t,
				test.backendOptions,
				nil,
				aladino.MockBuiltIns(),
				nil,
			)

			res, err := isMerged(mockedEnv, []aladino.Value{})

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantResult, res)
		})
	}
}
