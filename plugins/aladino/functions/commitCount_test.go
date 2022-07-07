// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commitCount = plugins_aladino.PluginBuiltIns().Functions["commitCount"].Code

func TestCommitCount(t *testing.T) {
	totalCommits := 1
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.Commits = &totalCommits
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		mock.WithRequestMatchHandler(
			// Overwrite default mock to pull request request details
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(mock.MustMarshal(defaultPullRequestDetails))
			}),
		),
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCommitCount := aladino.BuildIntValue(*mockedEnv.GetPullRequest().Commits)

	args := []aladino.Value{}
	gotCommitCount, err := commitCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCommitCount, gotCommitCount)
}
