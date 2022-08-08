// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	protobuf "github.com/explore-dev/atlas-common/go/api/services"
	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/shurcooL/githubv4"
)

type Severity int

const SEVERITY_FATAL Severity = 1

type TypeEnv map[string]Type

type Patch map[string]*File

type RegisterMap map[string]Value

type Env interface {
	GetBuiltIns() *BuiltIns
	GetBuiltInsReportedMessages() map[Severity][]string
	GetClient() *github.Client
	GetClientGQL() *githubv4.Client
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
	BuiltInsReportedMessages map[Severity][]string
	Ctx                      context.Context
	DryRun                   bool
	Client                   *github.Client
	ClientGQL                *githubv4.Client
	Collector                collector.Collector
	PullRequest              *github.PullRequest
	Patch                    Patch
	RegisterMap              RegisterMap
	BuiltIns                 *BuiltIns
	Report                   *Report
	EventPayload             interface{}
}

func (e *BaseEnv) GetBuiltIns() *BuiltIns {
	return e.BuiltIns
}

func (e *BaseEnv) GetBuiltInsReportedMessages() map[Severity][]string {
	return e.BuiltInsReportedMessages
}

func (e *BaseEnv) GetClient() *github.Client {
	return e.Client
}

func (e *BaseEnv) GetClientGQL() *githubv4.Client {
	return e.ClientGQL
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
	gitHubClient *github.Client,
	gitHubClientGQL *githubv4.Client,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	eventPayload interface{},
	builtIns *BuiltIns,
	semanticClient protobuf.SemanticClient,
) (Env, error) {
	owner := utils.GetPullRequestBaseOwnerName(pullRequest)
	repo := utils.GetPullRequestBaseRepoName(pullRequest)
	number := utils.GetPullRequestNumber(pullRequest)

	files, err := utils.GetPullRequestFiles(ctx, gitHubClient, owner, repo, number)
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
	}

	input := &BaseEnv{
		BuiltInsReportedMessages: builtInsReportedMessages,
		Ctx:                      ctx,
		DryRun:                   dryRun,
		Client:                   gitHubClient,
		ClientGQL:                gitHubClientGQL,
		Collector:                collector,
		PullRequest:              pullRequest,
		Patch:                    patch,
		RegisterMap:              registerMap,
		BuiltIns:                 builtIns,
		Report:                   report,
		EventPayload:             eventPayload,
	}

	return input, nil
}
