// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewEvalEnv(t *testing.T) {
	ctx := engine.DefaultMockCtx
	client := engine.MockGithubClient(nil)
	collector := engine.DefaultMockCollector
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()
	fileName := "default-mock-repo/file1.ts"
	patch := `@@ -2,9 +2,11 @@ package main
- func previous1() {
+ func new1() {
+
return`

	mockedFile := &aladino.File{
		Repr: &github.CommitFile{
			Filename: github.String(fileName),
			Patch:    github.String(patch),
		},
	}
	mockedFile.AppendToDiff(false, 2, 2, 2, 3, " func previous1() {", " func new1() {\n")
	mockedPatch := aladino.Patch{
		fileName: mockedFile,
	}

	mockedAladinoInterpreter := &aladino.Interpreter{
		Env: &aladino.BaseEnv{
			Ctx:          ctx,
			Client:       client,
			ClientGQL:    nil,
			Collector:    collector,
			PullRequest:  mockedPullRequest,
			Patch:        mockedPatch,
			RegisterMap:  aladino.RegisterMap(make(map[string]aladino.Value)),
			BuiltIns:     aladino.MockBuiltIns(),
			Report:       &aladino.Report{WorkflowDetails: make(map[string]aladino.ReportWorkflowDetails)},
			EventPayload: engine.DefaultMockEventPayload,
		},
	}

	wantEnv := &engine.Env{
		Ctx:          ctx,
		DryRun:       aladino.DefaultMockDryRun,
		Client:       client,
		ClientGQL:    nil,
		Collector:    collector,
		PullRequest:  mockedPullRequest,
		EventPayload: engine.DefaultMockEventPayload,
		Interpreter:  mockedAladinoInterpreter,
	}

	gotEnv, err := engine.NewEvalEnv(
		ctx,
		aladino.DefaultMockDryRun,
		client,
		nil,
		collector,
		mockedPullRequest,
		engine.DefaultMockEventPayload,
		mockedAladinoInterpreter,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantEnv, gotEnv)
}
