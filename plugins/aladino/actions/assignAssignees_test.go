// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/stretchr/testify/assert"
)

var assignAssignees = plugins_aladino.PluginBuiltIns().Actions["assignAssignees"].Code

type AssigneesRequestPostBody struct {
	Assignees []string `json:"assignees"`
}

func TestAssignAssignees_WhenListOfAssigneesIsEmpty(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err := assignAssignees(mockedEnv, args)

	assert.EqualError(t, err, "assignAssignees: list of assignees can't be empty")
}

func TestAssignAssignees_WhenListOfAssigneesExceeds10Users(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("john"),
		aladino.BuildStringValue("mary"),
		aladino.BuildStringValue("jane"),
		aladino.BuildStringValue("steve"),
		aladino.BuildStringValue("peter"),
		aladino.BuildStringValue("adam"),
		aladino.BuildStringValue("nancy"),
		aladino.BuildStringValue("susan"),
		aladino.BuildStringValue("bob"),
		aladino.BuildStringValue("michael"),
		aladino.BuildStringValue("tom"),
	})}
	err := assignAssignees(mockedEnv, args)

	assert.EqualError(t, err, "assignAssignees: can only assign up to 10 assignees")
}

func TestAssignAssignees(t *testing.T) {
	gotAssignees := []string{}
	wantAssignees := []string{
		"mary",
	}
	mockedPullRequest := aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Assignees: []*github.User{},
	})

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesAssigneesByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := AssigneesRequestPostBody{}

					utils.MustUnmarshal(rawBody, &body)

					gotAssignees = body.Assignees
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		nil,
	)

	assignees := make([]aladino.Value, len(wantAssignees))
	for i, assignee := range wantAssignees {
		assignees[i] = aladino.BuildStringValue(assignee)
	}

	args := []aladino.Value{aladino.BuildArrayValue(assignees)}
	err := assignAssignees(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantAssignees, gotAssignees)
}
