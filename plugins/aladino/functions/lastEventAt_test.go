// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"
	"time"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var lastEventAt = plugins_aladino.PluginBuiltIns().Functions["lastEventAt"].Code

func TestLastEventAt(t *testing.T) {
	date := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	mockedCodeReview := aladino.GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		UpdatedAt: timestamppb.New(date),
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

	wantUpdatedAt := lang.BuildIntValue(int(date.Unix()))

	args := []lang.Value{}
	gotUpdatedAt, gotErr := lastEventAt(mockedEnv, args)

	assert.Nil(t, gotErr)
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}
