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

var title = plugins_aladino.PluginBuiltIns().Functions["title"].Code

func TestTitle(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	prTitle := "Amazing new feature"
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Title: github.String(prTitle),
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
		aladino.DefaultMockEventPayload,
		controller,
	)

	wantTitle := aladino.BuildStringValue(prTitle)

	args := []aladino.Value{}
	gotTitle, err := title(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTitle, gotTitle)
}
