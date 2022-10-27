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
	"github.com/stretchr/testify/assert"
)

var labels = plugins_aladino.PluginBuiltIns().Functions["labels"].Code

func TestLabels(t *testing.T) {
	ghLabels := []*github.Label{
		{
			Name: github.String("bug"),
			ID:   github.Int64(1),
		},
		{
			Name: github.String("enhancement"),
			ID:   github.Int64(2),
		},
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Labels: ghLabels,
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

	wantLabels := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("bug"),
		aladino.BuildStringValue("enhancement"),
	})

	args := []aladino.Value{}
	gotLabels, err := labels(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}
