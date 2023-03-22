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

var head = plugins_aladino.PluginBuiltIns().Functions["head"].Code

func TestHead(t *testing.T) {
	headRef := "new-topic"
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Head: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: "john",
				Name:  "default-mock-repo",
			},
			Name: headRef,
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

	wantHead := aladino.BuildStringValue(headRef)

	args := []aladino.Value{}
	gotHead, err := head(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantHead, gotHead, "it should get the pull request head reference")
}
