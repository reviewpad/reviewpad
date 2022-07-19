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
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v3/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var size = plugins_aladino.PluginBuiltIns().Functions["size"].Code

func TestSize(t *testing.T) {
	additions := 1
	deletions := 1
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Additions: github.Int(additions),
		Deletions: github.Int(deletions),
	})
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantSize := aladino.BuildIntValue(additions + deletions)

	args := []aladino.Value{}
	gotSize, err := size(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantSize, gotSize)
}
