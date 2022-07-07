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

var base = plugins_aladino.PluginBuiltIns().Functions["base"].Code

func TestBase(t *testing.T) {
  ownerLogin := "john"
  repoName := "default-mock-repo"
	baseRef := "master"
	defaultPullRequestDetails := mocks_aladino.GetDefaultMockPullRequestDetails()
	defaultPullRequestDetails.Base.Repo.Owner = &github.User{
        Login: github.String(ownerLogin),
  }
  defaultPullRequestDetails.Base.Repo.Name = github.String(repoName)
  defaultPullRequestDetails.Base.Ref = github.String(baseRef)
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

	wantBase := aladino.BuildStringValue(mockedEnv.GetPullRequest().GetBase().GetRef())

	args := []aladino.Value{}
	gotBase, err := base(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantBase, gotBase)
}
