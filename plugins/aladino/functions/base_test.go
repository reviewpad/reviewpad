// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var base = plugins_aladino.PluginBuiltIns().Functions["base"].Code

func TestBase(t *testing.T) {
	baseRef := "master"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: "john",
				Name:  "default-mock-repo",
			},
			Name: baseRef,
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

	wantBase := aladino.BuildStringValue(baseRef)

	args := []aladino.Value{}
	gotBase, err := base(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantBase, gotBase, "it should get the pull request base reference")
}
