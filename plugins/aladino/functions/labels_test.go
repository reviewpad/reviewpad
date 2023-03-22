// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var labels = plugins_aladino.PluginBuiltIns().Functions["labels"].Code

func TestLabels(t *testing.T) {
	ghLabels := []*pbe.Label{
		{
			Name: "bug",
			Id:   "1",
		},
		{
			Name: "enhancement",
			Id:   "2",
		},
	}
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Labels: ghLabels,
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantLabels := aladino.BuildArrayValue([]aladino.Value{
		aladino.BuildStringValue("bug"),
		aladino.BuildStringValue("enhancement"),
	})

	args := []aladino.Value{}
	gotLabels, err := labels(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLabels, gotLabels)
}
