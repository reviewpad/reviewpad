// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var eventType = plugins_aladino.PluginBuiltIns().Functions["eventType"].Code

func TestEventType_WhenEventPayloadIsNotPullRequestEvent(t *testing.T) {
	wantValue := lang.BuildStringValue("")

	eventPayload := &github.CheckRunEvent{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{}
	gotValue, err := eventType(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestEventType_WhenPullRequestEventIsNil(t *testing.T) {
	wantValue := lang.BuildStringValue("")

	eventPayload := &github.PullRequestEvent{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{}
	gotValue, err := eventType(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestEventType_WhenPullRequestEventActionIsNil(t *testing.T) {
	wantValue := lang.BuildStringValue("")

	eventPayload := &github.PullRequestEvent{
		Action: nil,
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{}
	gotValue, err := eventType(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestEventType(t *testing.T) {
	wantValue := lang.BuildStringValue("synchronized")

	eventPayload := &github.PullRequestEvent{
		Action: github.String("synchronized"),
	}

	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{}
	gotValue, err := eventType(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}
