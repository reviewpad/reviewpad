// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignAssignees = plugins_aladino.PluginBuiltIns().Actions["assignAssignees"].Code

type AssigneesRequestPostBody struct {
	Assignees []string `json:"assignees"`
}

func TestAssignAssignees_WhenListOfAssigneesIsEmpty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err := assignAssignees(mockedEnv, args)

	assert.EqualError(t, err, "assignAssignees: list of assignees can't be empty")
}

func TestAssignAssignees_WhenListOfAssigneesExceeds10Users(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

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
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesAssigneesByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := AssigneesRequestPostBody{}

					json.Unmarshal(rawBody, &body)

					gotAssignees = body.Assignees
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		aladino.DefaultMockEventPayload,
		controller,
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
