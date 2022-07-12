// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commits = plugins_aladino.PluginBuiltIns().Functions["commits"].Code

func TestCommits_WhenListCommitsRequestFails(t *testing.T) {
	failMessage := "ListCommitsRequestFail"
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	gotCommits, err := commits(mockedEnv, args)

	assert.Nil(t, gotCommits)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommits(t *testing.T) {
	repoCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
			},
		},
	}
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				repoCommits,
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantCommitsMessages := make([]aladino.Value, len(repoCommits))
	for i, repoCommit := range repoCommits {
		wantCommitsMessages[i] = aladino.BuildStringValue(repoCommit.Commit.GetMessage())
	}
	wantCommits := aladino.BuildArrayValue(wantCommitsMessages)

	args := []aladino.Value{}
	gotCommits, err := commits(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCommits, gotCommits)
}
