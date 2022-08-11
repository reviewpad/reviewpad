// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var hasLinearHistory = plugins_aladino.PluginBuiltIns(plugins_aladino.DefaultPluginConfig()).Functions["hasLinearHistory"].Code

func TestHasLinearHistory_WhenListCommitsRequestFails(t *testing.T) {
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
		nil,
	)

	args := []aladino.Value{}
	gotVal, err := hasLinearHistory(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestHasLinearHistory_WhenFalse(t *testing.T) {
	repoCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
				Parents: []*github.Commit{
					{Message: github.String("FirstParent")},
					{Message: github.String("SecondParent")},
				},
			},
			Parents: []*github.Commit{
				{Message: github.String("FirstParent")},
				{Message: github.String("SecondParent")},
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
		nil,
	)

	args := []aladino.Value{}
	gotVal, err := hasLinearHistory(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestHasLinearHistory_WhenTrue(t *testing.T) {
	repoCommits := []*github.RepositoryCommit{
		{
			Commit: &github.Commit{
				Message: github.String("Lorem Ipsum"),
				Parents: []*github.Commit{
					{Message: github.String("OnlyParent")},
				},
			},
			Parents: []*github.Commit{
				{Message: github.String("OnlyParent")},
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
		nil,
	)

	args := []aladino.Value{}
	gotVal, err := hasLinearHistory(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
