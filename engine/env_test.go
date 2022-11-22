// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewEvalEnv(t *testing.T) {
	ctx := engine.DefaultMockCtx
	githubClient := engine.MockGithubClient(nil)
	collector := engine.DefaultMockCollector

	mockedAladinoInterpreter := &aladino.Interpreter{
		Env: &aladino.BaseEnv{
			Ctx:          ctx,
			GithubClient: githubClient,
			Collector:    collector,
			RegisterMap:  aladino.RegisterMap(make(map[string]aladino.Value)),
			BuiltIns:     aladino.MockBuiltIns(),
			Report:       &aladino.Report{Actions: make([]string, 0)},
		},
	}

	wantEnv := &engine.Env{
		Ctx:          ctx,
		DryRun:       false,
		GithubClient: githubClient,
		Collector:    collector,
		Interpreter:  mockedAladinoInterpreter,
		TargetEntity: aladino.DefaultMockTargetEntity,
		EventData:    aladino.DefaultMockEventData,
	}

	gotEnv, err := engine.NewEvalEnv(
		ctx,
		false,
		githubClient,
		collector,
		aladino.DefaultMockTargetEntity,
		mockedAladinoInterpreter,
		aladino.DefaultMockEventData,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantEnv, gotEnv)
}
