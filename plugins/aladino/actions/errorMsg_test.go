// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var errorMsg = plugins_aladino.PluginBuiltIns().Actions["error"].Code

func TestErrorMsg(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	message := "This pull request was considered too large."

	args := []aladino.Value{aladino.BuildStringValue(message)}
	gotError := errorMsg(mockedEnv, args)

	wantReportedErrors := []string{message}

	gotReportedErrors := mockedEnv.GetBuiltInsReportedMessages()[aladino.SEVERITY_ERROR]

	assert.Nil(t, gotError)
	assert.Equal(t, wantReportedErrors, gotReportedErrors)
}
