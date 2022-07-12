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

var assignees = plugins_aladino.PluginBuiltIns().Functions["assignees"].Code

func TestAssignees(t *testing.T) {
	assigneeLogin := "jane"
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Assignees: []*github.User{
			{Login: github.String(assigneeLogin)},
		},
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

	mockedAssignees := mockedEnv.GetPullRequest().Assignees
	wantAssigneesLogins := make([]aladino.Value, len(mockedAssignees))
	for i, assignee := range mockedAssignees {
		wantAssigneesLogins[i] = aladino.BuildStringValue(assignee.GetLogin())
	}

	wantAssignees := aladino.BuildArrayValue(wantAssigneesLogins)

	gotAssignees, err := assignees(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantAssignees, gotAssignees)
}
