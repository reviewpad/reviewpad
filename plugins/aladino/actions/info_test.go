// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var info = plugins_aladino.PluginBuiltIns().Actions["info"].Code

func TestInfo(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil, handler.PullRequest)

	message := "This pull request was considered too large."

	args := []aladino.Value{aladino.BuildStringValue(message)}
	gotError := info(mockedEnv, args)

	wantReportedInfos := []string{message}

	gotReportedInfos := mockedEnv.GetBuiltInsReportedMessages()[aladino.SEVERITY_INFO]

	assert.Nil(t, gotError)
	assert.Equal(t, wantReportedInfos, gotReportedInfos)
}
