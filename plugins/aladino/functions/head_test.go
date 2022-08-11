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

var head = plugins_aladino.PluginBuiltIns(plugins_aladino.DefaultPluginConfig()).Functions["head"].Code

func TestHead(t *testing.T) {
	headRef := "new-topic"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("john"),
				},
				Name: github.String("default-mock-repo"),
			},
			Ref: github.String(headRef),
		},
	})
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	wantHead := aladino.BuildStringValue(headRef)

	args := []aladino.Value{}
	gotHead, err := head(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantHead, gotHead, "it should get the pull request head reference")
}
