// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestNewEvalEnv(t *testing.T) {
	ctx := engine.DefaultMockCtx
	githubClient := engine.MockGithubClient(nil)
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
			GithubClient: githubClient,
			Collector:    collector,
			PullRequest:  mockedPullRequest,
			Patch:        mockedPatch,
			RegisterMap:  aladino.RegisterMap(make(map[string]aladino.Value)),
			BuiltIns:     aladino.MockBuiltIns(),
			Report:       &aladino.Report{Actions: make([]string, 0)},
			EventPayload: engine.DefaultMockEventPayload,
		},
	}

	wantEnv := &engine.Env{
		Ctx:          ctx,
		DryRun:       false,
		GithubClient: githubClient,
		Collector:    collector,
		PullRequest:  mockedPullRequest,
		EventPayload: engine.DefaultMockEventPayload,
		Interpreter:  mockedAladinoInterpreter,
		TargetEntity: aladino.DefaultMockTargetEntity,
	}

	gotEnv, err := engine.NewEvalEnv(
		ctx,
		false,
		githubClient,
		collector,
		engine.DefaultMockTargetEntity,
		engine.DefaultMockEventPayload,
		mockedAladinoInterpreter,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantEnv, gotEnv)
}

func TestNewEvalEnv_WhenGetPullRequestFails(t *testing.T) {
	failMessage := "GetPullRequestFail"
	ctx := engine.DefaultMockCtx
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
	mockedGithubClientREST := github.NewClient(mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposPullsByOwnerByRepoByPullNumber,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mock.WriteError(
					w,
					http.StatusInternalServerError,
					failMessage,
				)
			}),
		),
	))

	mockedGithubClient := gh.NewGithubClient(mockedGithubClientREST, nil)

	mockedAladinoInterpreter := &aladino.Interpreter{
		Env: &aladino.BaseEnv{
			Ctx:          ctx,
			GithubClient: mockedGithubClient,
			Collector:    engine.DefaultMockCollector,
			PullRequest:  engine.GetDefaultMockPullRequestDetails(),
			Patch:        mockedPatch,
			RegisterMap:  aladino.RegisterMap(make(map[string]aladino.Value)),
			BuiltIns:     aladino.MockBuiltIns(),
			Report:       &aladino.Report{Actions: make([]string, 0)},
			EventPayload: engine.DefaultMockEventPayload,
		},
	}

	env, err := engine.NewEvalEnv(
		ctx,
		false,
		mockedGithubClient,
		engine.DefaultMockCollector,
		engine.DefaultMockTargetEntity,
		engine.DefaultMockTargetEntity,
		mockedAladinoInterpreter,
	)

	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
	assert.Nil(t, env)
}
