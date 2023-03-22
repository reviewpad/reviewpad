// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"
	"time"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var lastEventAt = plugins_aladino.PluginBuiltIns().Functions["lastEventAt"].Code

func TestLastEventAt(t *testing.T) {
	date := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		UpdatedAt: timestamppb.New(date),
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantUpdatedAt := aladino.BuildIntValue(int(date.Unix()))

	args := []aladino.Value{}
	gotUpdatedAt, gotErr := lastEventAt(mockedEnv, args)

	assert.Nil(t, gotErr)
	assert.Equal(t, wantUpdatedAt, gotUpdatedAt)
}
