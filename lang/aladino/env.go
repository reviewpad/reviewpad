// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/lang/aladino/target"
)

type Severity int

const (
	SEVERITY_FATAL   Severity = 1
	SEVERITY_ERROR   Severity = 2
	SEVERITY_WARNING Severity = 3
	SEVERITY_INFO    Severity = 4
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
	GetIssue() *github.Issue
	GetTargetEntity() *handler.TargetEntity
	GetTarget() target.Target
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
	TargetEntity             *handler.TargetEntity
	Issue                    *github.Issue
	Target                   target.Target
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

func (e *BaseEnv) GetTargetEntity() *handler.TargetEntity {
	return e.TargetEntity
}

func (e *BaseEnv) GetIssue() *github.Issue {
	return e.Issue
}

func (e *BaseEnv) GetTarget() target.Target {
	return e.Target
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

func getPullRequestPatch(ctx context.Context, pullRequest *github.PullRequest, githubClient *gh.GithubClient) (map[string]*File, error) {
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

	return Patch(patchMap), nil
}

func NewEvalEnv(
	ctx context.Context,
	dryRun bool,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	targetEntity *handler.TargetEntity,
	eventPayload interface{},
	builtIns *BuiltIns,
) (Env, error) {
	registerMap := RegisterMap(make(map[string]Value))
	report := &Report{Actions: make([]string, 0)}

	input := &BaseEnv{
		BuiltIns:                 builtIns,
		BuiltInsReportedMessages: make(map[Severity][]string),
		GithubClient:             githubClient,
		Collector:                collector,
		Ctx:                      ctx,
		DryRun:                   dryRun,
		EventPayload:             eventPayload,
		RegisterMap:              registerMap,
		Report:                   report,
		TargetEntity:             targetEntity,
	}

	if targetEntity.Kind == handler.Issue {
		issue, _, err := githubClient.GetIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		if err != nil {
			return nil, err
		}

		input.Issue = issue
		input.Target = target.NewIssueTarget(ctx, targetEntity, githubClient, issue)

		return input, nil
	}

	pullRequest, _, err := githubClient.GetPullRequest(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
	if err != nil {
		return nil, err
	}

	patch, err := getPullRequestPatch(ctx, pullRequest, githubClient)
	if err != nil {
		return nil, err
	}

	input.Patch = patch
	input.PullRequest = pullRequest
	input.Target = target.NewPullRequestTarget(ctx, targetEntity, githubClient, pullRequest)

	return input, nil
}
