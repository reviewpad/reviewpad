// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/shurcooL/githubv4"
)

type TypeEnv map[string]Type

type Patch map[string]*File

type RegisterMap map[string]Value

type Env interface {
	GetCtx() context.Context
	GetClient() *github.Client
	GetClientGQL() *githubv4.Client
	GetCollector() collector.Collector
	GetPullRequest() *github.PullRequest
	GetPatch() Patch
	GetRegisterMap() RegisterMap
	GetBuiltIns() *BuiltIns
	GetReport() *Report
	GetEventPayload() interface{}
	GetDryRun() bool
}

type BaseEnv struct {
	Ctx          context.Context
	DryRun       bool
	Client       *github.Client
	ClientGQL    *githubv4.Client
	Collector    collector.Collector
	PullRequest  *github.PullRequest
	Patch        Patch
	RegisterMap  RegisterMap
	BuiltIns     *BuiltIns
	Report       *Report
	EventPayload interface{}
}

func (e *BaseEnv) GetCtx() context.Context {
	return e.Ctx
}

func (e *BaseEnv) GetDryRun() bool {
	return e.DryRun
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

func (e *BaseEnv) GetPullRequest() *github.PullRequest {
	return e.PullRequest
}

func (e *BaseEnv) GetPatch() Patch {
	return e.Patch
}

func (e *BaseEnv) GetRegisterMap() RegisterMap {
	return e.RegisterMap
}

func (e *BaseEnv) GetBuiltIns() *BuiltIns {
	return e.BuiltIns
}

func (e *BaseEnv) GetReport() *Report {
	return e.Report
}

func (e *BaseEnv) GetEventPayload() interface{} {
	return e.EventPayload
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
	report := &Report{WorkflowDetails: make(map[string]ReportWorkflowDetails, 0)}

	input := &BaseEnv{
		Ctx:          ctx,
		DryRun:       dryRun,
		Client:       gitHubClient,
		ClientGQL:    gitHubClientGQL,
		Collector:    collector,
		PullRequest:  pullRequest,
		Patch:        patch,
		RegisterMap:  registerMap,
		BuiltIns:     builtIns,
		Report:       report,
		EventPayload: eventPayload,
	}

	return input, nil
}
