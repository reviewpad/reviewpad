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

var state = plugins_aladino.PluginBuiltIns().Functions["state"].Code

func TestStatus(t *testing.T) {
	prState := "open"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Status: pbe.CodeReviewStatus_OPEN,
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantState := aladino.BuildStringValue(prState)

	args := []aladino.Value{}
	gotState, err := state(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantState, gotState)
}
