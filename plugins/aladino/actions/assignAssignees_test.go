// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"encoding/json"
	"io/ioutil"
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

var assignAssignees = plugins_aladino.PluginBuiltIns().Actions["assignAssignees"].Code

type AssigneesRequestPostBody struct {
	Assignees []string `json:"assignees"`
}

func TestAssignAssignees_WhenListOfAssigneesIsEmpty(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{})}
	err = assignAssignees(mockedEnv, args)

	assert.EqualError(t, err, "assignAssignees: list of assignees can't be empty")
}

func TestAssignAssignees_WhenListOfAssigneesExceeds10Users(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

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
	err = assignAssignees(mockedEnv, args)

	assert.EqualError(t, err, "assignAssignees: can only assign up to 10 assignees")
}

func TestAssignAssignees(t *testing.T) {
	gotAssignees := []string{}
	wantAssignees := []string{
		"mary",
	}
	mockedPullRequest := mocks_aladino.GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		Assignees: []*github.User{},
	})

	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	assignees := make([]aladino.Value, len(wantAssignees))
	for i, assignee := range wantAssignees {
		assignees[i] = aladino.BuildStringValue(assignee)
	}

	args := []aladino.Value{aladino.BuildArrayValue(assignees)}
	err = assignAssignees(mockedEnv, args)

	assert.Nil(t, err)
	assert.ElementsMatch(t, wantAssignees, gotAssignees)
}
