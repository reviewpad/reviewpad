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

var createdAt = plugins_aladino.PluginBuiltIns().Functions["createdAt"].Code

func TestCreatedAt(t *testing.T) {
	date := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		CreatedAt: timestamppb.New(date),
	})
	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		nil,
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	wantCreatedAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", date.String())
	if err != nil {
		assert.FailNow(t, "time.Parse failed", err)
	}
	wantCreatedAt := aladino.BuildIntValue(int(wantCreatedAtTime.Unix()))

	args := []aladino.Value{}
	gotCreatedAt, err := createdAt(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCreatedAt, gotCreatedAt)
}
