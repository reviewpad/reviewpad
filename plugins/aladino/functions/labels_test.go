// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var labels = plugins_aladino.PluginBuiltIns().Functions["labels"].Code

func TestLabels(t *testing.T) {
	ghLabels := []*pbc.Label{
		{
			Name: "bug",
			Id:   "1",
		},
		{
			Name: "enhancement",
			Id:   "2",
		},
	}
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Labels: ghLabels,
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

	wantLabels := lang.BuildArrayValue([]lang.Value{
		lang.BuildStringValue("bug"),
		lang.BuildStringValue("enhancement"),
	})

	args := []lang.Value{}
	gotLabels, err := labels(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}
