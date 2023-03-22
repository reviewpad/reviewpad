// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var assignees = plugins_aladino.PluginBuiltIns().Functions["assignees"].Code

func TestAssignees(t *testing.T) {
	assigneeLogin := "jane"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Assignees: []*pbe.ExternalUser{
			{Login: assigneeLogin},
		},
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	mockedAssignees := mockedEnv.GetTarget().(*target.PullRequestTarget).CodeReview.Assignees
	wantAssigneesLogins := make([]aladino.Value, len(mockedAssignees))
	for i, assignee := range mockedAssignees {
		wantAssigneesLogins[i] = aladino.BuildStringValue(assignee.GetLogin())
	}

	wantAssignees := aladino.BuildArrayValue(wantAssigneesLogins)

	gotAssignees, err := assignees(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantAssignees, gotAssignees)
}
