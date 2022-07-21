// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewEvalEnv(t *testing.T) {
	ctx := engine.DefaultCtx
	client := engine.MockGithubClient(nil)
	collector := engine.DefaultCollector
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		ctx,
		client,
		nil,
		collector,
		mockedPullRequest,
		nil,
		// TODO: This needs to be mocked (PR #155)
		plugins_aladino.PluginBuiltIns(),
	)
	if err != nil {
		log.Fatalf("aladino NewInterpreter failed: %v", err)
	}

	wantEnv := &engine.Env{
		Ctx:          ctx,
		DryRun:       false,
		Client:       client,
		ClientGQL:    nil,
		Collector:    collector,
		PullRequest:  mockedPullRequest,
		EventPayload: nil,
		Interpreter:  mockedAladinoInterpreter,
	}

	gotEnv, err := engine.NewEvalEnv(
		ctx,
		false,
		client,
		nil,
		collector,
		mockedPullRequest,
		nil,
		mockedAladinoInterpreter,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantEnv, gotEnv)
}
