// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"

	"github.com/reviewpad/reviewpad/v3/codehost"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/sirupsen/logrus"
)

type Severity int

const (
	SEVERITY_FATAL   Severity = 1
	SEVERITY_ERROR   Severity = 2
	SEVERITY_WARNING Severity = 3
	SEVERITY_INFO    Severity = 4
)

type TypeEnv map[string]Type

type RegisterMap map[string]Value

type Env interface {
	GetBuiltIns() *BuiltIns
	GetBuiltInsReportedMessages() map[Severity][]string
	GetGithubClient() *gh.GithubClient
	GetCollector() collector.Collector
	GetCtx() context.Context
	GetDryRun() bool
	GetEventPayload() interface{}
	GetRegisterMap() RegisterMap
	GetReport() *Report
	GetTarget() codehost.Target
	GetLogger() *logrus.Entry
}

type BaseEnv struct {
	BuiltIns                 *BuiltIns
	BuiltInsReportedMessages map[Severity][]string
	GithubClient             *gh.GithubClient
	Collector                collector.Collector
	Ctx                      context.Context
	DryRun                   bool
	EventPayload             interface{}
	RegisterMap              RegisterMap
	Report                   *Report
	Target                   codehost.Target
	Logger                   *logrus.Entry
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

func (e *BaseEnv) GetRegisterMap() RegisterMap {
	return e.RegisterMap
}

func (e *BaseEnv) GetReport() *Report {
	return e.Report
}

func (e *BaseEnv) GetTarget() codehost.Target {
	return e.Target
}

func (e *BaseEnv) GetLogger() *logrus.Entry {
	return e.Logger
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
	logger *logrus.Entry,
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
		Logger:                   logger,
	}

	switch targetEntity.Kind {
	case handler.Issue:

		issue, _, err := githubClient.GetIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		if err != nil {
			return nil, err
		}

		input.Target = target.NewIssueTarget(ctx, targetEntity, githubClient, issue)
	case handler.PullRequest:
		pullRequest, _, err := githubClient.GetPullRequest(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		if err != nil {
			return nil, err
		}

		pullRequestTarget, err := target.NewPullRequestTarget(ctx, targetEntity, githubClient, pullRequest)
		if err != nil {
			return nil, err
		}
		input.Target = pullRequestTarget
	}

	return input, nil
}
