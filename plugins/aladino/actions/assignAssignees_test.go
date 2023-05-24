// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignAssignees = plugins_aladino.PluginBuiltIns().Actions["assignAssignees"].Code

func TestAssignAssignees(t *testing.T) {
	defaultAssignees := 10
	tests := map[string]struct {
		clientOptions               []mock.MockBackendOption
		codeReview                  *pbc.PullRequest
		inputAssignees              lang.Value
		inputTotalRequiredAssignees lang.Value
		shouldAssign                bool
		wantErr                     error
	}{
		"when list of assignees is empty": {
			inputAssignees:              lang.BuildArrayValue([]lang.Value{}),
			inputTotalRequiredAssignees: lang.BuildIntValue(defaultAssignees),
			wantErr:                     fmt.Errorf("assignAssignees: list of assignees can't be empty"),
		},
		"when list of assignees exceeds 10 users": {
			inputAssignees: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("john"),
				lang.BuildStringValue("mary"),
				lang.BuildStringValue("jane"),
				lang.BuildStringValue("steve"),
				lang.BuildStringValue("peter"),
				lang.BuildStringValue("adam"),
				lang.BuildStringValue("nancy"),
				lang.BuildStringValue("susan"),
				lang.BuildStringValue("bob"),
				lang.BuildStringValue("michael"),
				lang.BuildStringValue("tom"),
			}),
			inputTotalRequiredAssignees: lang.BuildIntValue(11),
			wantErr:                     fmt.Errorf("assignAssignees: can only assign up to 10 assignees"),
		},
		"when the total required assignees is an invalid number": {
			inputAssignees: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("john"),
				lang.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: lang.BuildIntValue(0),
			wantErr:                     fmt.Errorf("assignAssignees: total required assignees is invalid. please insert a number bigger than 0"),
		},
		"when the total required assignees is greater than the total of available assignees": {
			codeReview: &pbc.PullRequest{
				Assignees: []*pbc.User{},
			},
			inputAssignees: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("john"),
				lang.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: lang.BuildIntValue(3),
			shouldAssign:                true,
		},
		"when one of the required assignees is already an assignee": {
			codeReview: &pbc.PullRequest{
				Assignees: []*pbc.User{
					{
						Login: "john",
					},
				},
			},
			inputAssignees: lang.BuildArrayValue([]lang.Value{
				lang.BuildStringValue("john"),
				lang.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: lang.BuildIntValue(1),
			shouldAssign:                true,
		},
	}

	for _, test := range tests {
		isAssigneesRequestPerformed := false
		mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
			t,
			append(
				[]mock.MockBackendOption{
					mock.WithRequestMatchHandler(
						mock.PostReposIssuesAssigneesByOwnerByRepoByIssueNumber,
						http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							// If the assign request was performed then the given users were assigned to the pull request
							isAssigneesRequestPerformed = true
						}),
					),
				},
				test.clientOptions...,
			),
			nil,
			aladino.GetDefaultMockPullRequestDetailsWith(test.codeReview),
			aladino.GetDefaultPullRequestFileList(),
			aladino.MockBuiltIns(),
			nil,
		)

		args := []lang.Value{test.inputAssignees, test.inputTotalRequiredAssignees}
		gotErr := assignAssignees(mockedEnv, args)

		assert.Equal(t, test.shouldAssign, isAssigneesRequestPerformed)
		assert.Equal(t, test.wantErr, gotErr)
	}
}
