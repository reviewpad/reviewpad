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

var description = plugins_aladino.PluginBuiltIns().Functions["description"].Code

func TestDescription(t *testing.T) {
	prDescription := "Please pull these awesome changes in!"
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Description: prDescription,
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

	wantDescription := lang.BuildStringValue(prDescription)

	args := []lang.Value{}
	gotDescription, err := description(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantDescription, gotDescription)
}
