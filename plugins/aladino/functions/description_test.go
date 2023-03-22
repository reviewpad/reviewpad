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

var description = plugins_aladino.PluginBuiltIns().Functions["description"].Code

func TestDescription(t *testing.T) {
	prDescription := "Please pull these awesome changes in!"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Description: prDescription,
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantDescription := aladino.BuildStringValue(prDescription)

	args := []aladino.Value{}
	gotDescription, err := description(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantDescription, gotDescription)
}
