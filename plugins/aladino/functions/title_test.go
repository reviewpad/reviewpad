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

var title = plugins_aladino.PluginBuiltIns().Functions["title"].Code

func TestTitle(t *testing.T) {
	prTitle := "Amazing new feature"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Title: prTitle,
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantTitle := aladino.BuildStringValue(prTitle)

	args := []aladino.Value{}
	gotTitle, err := title(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantTitle, gotTitle)
}
