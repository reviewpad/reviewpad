// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignAssignees = plugins_aladino.PluginBuiltIns().Actions["assignAssignees"].Code

func TestAssignAssignees(t *testing.T) {
	defaultAssignees := 10
	tests := map[string]struct {
		clientOptions               []mock.MockBackendOption
		codeReview                  *pbe.CodeReview
		inputAssignees              aladino.Value
		inputTotalRequiredAssignees aladino.Value
		shouldAssign                bool
		wantErr                     error
	}{
		"when list of assignees is empty": {
			inputAssignees:              aladino.BuildArrayValue([]aladino.Value{}),
			inputTotalRequiredAssignees: aladino.BuildIntValue(defaultAssignees),
			wantErr:                     fmt.Errorf("assignAssignees: list of assignees can't be empty"),
		},
		"when list of assignees exceeds 10 users": {
			inputAssignees: aladino.BuildArrayValue([]aladino.Value{
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
			}),
			inputTotalRequiredAssignees: aladino.BuildIntValue(11),
			wantErr:                     fmt.Errorf("assignAssignees: can only assign up to 10 assignees"),
		},
		"when the total required assignees is an invalid number": {
			inputAssignees: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("john"),
				aladino.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: aladino.BuildIntValue(0),
			wantErr:                     fmt.Errorf("assignAssignees: total required assignees is invalid. please insert a number bigger than 0"),
		},
		"when the total required assignees is greater than the total of available assignees": {
			codeReview: &pbe.CodeReview{
				Assignees: []*pbe.ExternalUser{},
			},
			inputAssignees: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("john"),
				aladino.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: aladino.BuildIntValue(3),
			shouldAssign:                true,
		},
		"when one of the required assignees is already an assignee": {
			codeReview: &pbe.CodeReview{
				Assignees: []*pbe.ExternalUser{
					{
						Login: "john",
					},
				},
			},
			inputAssignees: aladino.BuildArrayValue([]aladino.Value{
				aladino.BuildStringValue("john"),
				aladino.BuildStringValue("mary"),
			}),
			inputTotalRequiredAssignees: aladino.BuildIntValue(1),
			shouldAssign:                true,
		},
	}

	for _, test := range tests {
		isAssigneesRequestPerformed := false
		mockedEnv := aladino.MockDefaultEnvWithCodeReview(
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
			aladino.GetDefaultMockCodeReviewDetailsWith(test.codeReview),
			aladino.MockBuiltIns(),
			nil,
		)

		args := []aladino.Value{test.inputAssignees, test.inputTotalRequiredAssignees}
		gotErr := assignAssignees(mockedEnv, args)

		assert.Equal(t, test.shouldAssign, isAssigneesRequestPerformed)
		assert.Equal(t, test.wantErr, gotErr)
	}
}
