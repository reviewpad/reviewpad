// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	"github.com/google/go-github/v45/github"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
)

type Severity int

const (
	SEVERITY_FATAL Severity = 1
	SEVERITY_ERROR Severity = 2
)

type TypeEnv map[string]Type

type Patch map[string]*File

type RegisterMap map[string]Value

type Env interface {
	GetBuiltIns() *BuiltIns
	GetBuiltInsReportedMessages() map[Severity][]string
	GetGithubClient() *gh.GithubClient
	GetCollector() collector.Collector
	GetCtx() context.Context
	GetDryRun() bool
	GetEventPayload() interface{}
	GetPatch() Patch
	GetPullRequest() *github.PullRequest
	GetRegisterMap() RegisterMap
	GetReport() *Report
}

type BaseEnv struct {
	BuiltIns                 *BuiltIns
	BuiltInsReportedMessages map[Severity][]string
	GithubClient             *gh.GithubClient
	Collector                collector.Collector
	Ctx                      context.Context
	DryRun                   bool
	EventPayload             interface{}
	Patch                    Patch
	PullRequest              *github.PullRequest
	RegisterMap              RegisterMap
	Report                   *Report
}

func (e *BaseEnv) GetBuiltIns() *BuiltIns {
	return e.BuiltIns
}

func (e *BaseEnv) GetBuiltInsReportedMessages() map[Severity][]string {
	return e.BuiltInsReportedMessages
}

func (e *BaseEnv) GetGithubClient() *gh.GithubClient {
	return e.GithubClient
}

func (e *BaseEnv) GetCollector() collector.Collector {
	return e.Collector
}

func (e *BaseEnv) GetCtx() context.Context {
	return e.Ctx
}

func (e *BaseEnv) GetDryRun() bool {
	return e.DryRun
}

func (e *BaseEnv) GetEventPayload() interface{} {
	return e.EventPayload
}

func (e *BaseEnv) GetPatch() Patch {
	return e.Patch
}

func (e *BaseEnv) GetPullRequest() *github.PullRequest {
	return e.PullRequest
}

func (e *BaseEnv) GetRegisterMap() RegisterMap {
	return e.RegisterMap
}

func (e *BaseEnv) GetReport() *Report {
	return e.Report
}

func NewTypeEnv(e Env) TypeEnv {
	builtInsType := make(map[string]Type)
	for builtInName, builtInFunction := range e.GetBuiltIns().Functions {
		builtInsType[builtInName] = builtInFunction.Type
	}

	for builtInName, builtInAction := range e.GetBuiltIns().Actions {
		builtInsType[builtInName] = builtInAction.Type
	}

	return TypeEnv(builtInsType)
}

func NewEvalEnv(
	ctx context.Context,
	dryRun bool,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	eventPayload interface{},
	builtIns *BuiltIns,
) (Env, error) {
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)
	number := gh.GetPullRequestNumber(pullRequest)

	files, err := githubClient.GetPullRequestFiles(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	patchMap := make(map[string]*File)

	for _, file := range files {
		patchFile, err := NewFile(file)
		if err != nil {
			return nil, err
		}

		patchMap[file.GetFilename()] = patchFile
	}

	patch := Patch(patchMap)
	registerMap := RegisterMap(make(map[string]Value))
	report := &Report{Actions: make([]string, 0)}
	builtInsReportedMessages := map[Severity][]string{
		SEVERITY_FATAL: make([]string, 0),
		SEVERITY_ERROR: make([]string, 0),
	}

	input := &BaseEnv{
		BuiltIns:                 builtIns,
		BuiltInsReportedMessages: builtInsReportedMessages,
		GithubClient:             githubClient,
		Collector:                collector,
		Ctx:                      ctx,
		DryRun:                   dryRun,
		EventPayload:             eventPayload,
		Patch:                    patch,
		PullRequest:              pullRequest,
		RegisterMap:              registerMap,
		Report:                   report,
	}

	return input, nil
}
