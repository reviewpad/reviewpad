// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commits = plugins_aladino.PluginBuiltIns().Functions["commits"].Code

func TestCommits_WhenListCommitsRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	failMessage := "ListCommitsRequestFail"
	mockedEnv := aladino.MockDefaultEnv(
		t,
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
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

	args := []aladino.Value{}
	gotCommits, err := commits(mockedEnv, args)

	assert.Nil(t, gotCommits)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestCommits(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repoCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposPullsCommitsByOwnerByRepoByPullNumber,
				repoCommits,
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
	)

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
