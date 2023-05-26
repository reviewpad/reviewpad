// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignees = plugins_aladino.PluginBuiltIns().Functions["assignees"].Code

func TestAssignees(t *testing.T) {
	assigneeLogin := "jane"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Assignees: []*pbc.User{
			{Login: assigneeLogin},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.GetDefaultPullRequestFileList(),
		aladino.MockBuiltIns(),
		nil,
	)

	mockedAssignees := mockedEnv.GetTarget().(*target.PullRequestTarget).PullRequest.Assignees
	wantAssigneesLogins := make([]lang.Value, len(mockedAssignees))
	for i, assignee := range mockedAssignees {
		wantAssigneesLogins[i] = lang.BuildStringValue(assignee.GetLogin())
	}

	wantAssignees := lang.BuildArrayValue(wantAssigneesLogins)

	gotAssignees, err := assignees(mockedEnv, []lang.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantAssignees, gotAssignees)
}
