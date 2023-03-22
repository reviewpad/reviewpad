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

var author = plugins_aladino.PluginBuiltIns().Functions["author"].Code

func TestAuthor(t *testing.T) {
	authorLogin := "john"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Author: &pbe.ExternalUser{Login: authorLogin},
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantAuthor := aladino.BuildStringValue(authorLogin)

	gotAuthor, err := author(mockedEnv, []aladino.Value{})

	assert.Nil(t, err)
	assert.Equal(t, wantAuthor, gotAuthor, "it should get the pull request author")
}
