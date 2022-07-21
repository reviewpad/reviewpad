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
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var labels = plugins_aladino.PluginBuiltIns().Functions["labels"].Code

func TestLabels(t *testing.T) {
	ghLabels := []*github.Label{
		{Name: github.String("bug")},
		{Name: github.String("enhancement")},
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Labels: ghLabels,
	})
	mockedEnv, err := aladino.MockDefaultEnv(
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

	wantLabels := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("bug"),
		aladino.BuildStringValue("enhancement"),
	})

	args := []aladino.Value{}
	gotLabels, err := labels(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}
