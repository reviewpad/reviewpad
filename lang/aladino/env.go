// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"sync"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/sirupsen/logrus"
)

type Severity int

const (
	SEVERITY_FAIL    Severity = 0
	SEVERITY_FATAL   Severity = 1
	SEVERITY_ERROR   Severity = 2
	SEVERITY_WARNING Severity = 3
	SEVERITY_INFO    Severity = 4
)

type TypeEnv map[string]lang.Type

type RegisterMap map[string]lang.Value

type Env interface {
	GetBuiltIns() *BuiltIns
	GetBuiltInsReportedMessages() map[Severity][]string
	GetGithubClient() *gh.GithubClient
	GetCodeHostClient() *codehost.CodeHostClient
	GetCollector() collector.Collector
	GetCtx() context.Context
	GetDryRun() bool
	GetEventPayload() interface{}
	GetRegisterMap() RegisterMap
	GetReport() *Report
	GetTarget() codehost.Target
	GetLogger() *logrus.Entry
	GetExecWaitGroup() *sync.WaitGroup
	GetExecFatalErrorOccurred() error
	SetExecFatalErrorOccurred(error)
	GetCheckRunID() *int64
	SetCheckRunConclusion(string)
	GetCheckRunConclusion() string
}

type BaseEnv struct {
	BuiltIns                 *BuiltIns
	BuiltInsReportedMessages map[Severity][]string
	GithubClient             *gh.GithubClient
	CodeHostClient           *codehost.CodeHostClient
	Collector                collector.Collector
	Ctx                      context.Context
	DryRun                   bool
	EventPayload             interface{}
	RegisterMap              RegisterMap
	Report                   *Report
	Target                   codehost.Target
	Logger                   *logrus.Entry
	ExecWaitGroup            *sync.WaitGroup
	ExecMutex                *sync.Mutex
	ExecFatalErrorOccurred   error
	CheckRunID               *int64
	CheckRunConclusion       string
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

func (e *BaseEnv) GetCodeHostClient() *codehost.CodeHostClient {
	return e.CodeHostClient
}

func (e *BaseEnv) GetExecWaitGroup() *sync.WaitGroup {
	return e.ExecWaitGroup
}

func (e *BaseEnv) SetCheckRunConclusion(conclusion string) {
	e.CheckRunConclusion = conclusion
}

func (e *BaseEnv) GetCheckRunConclusion() string {
	return e.CheckRunConclusion
}

func (e *BaseEnv) GetExecFatalErrorOccurred() error {
	e.ExecMutex.Lock()
	defer e.ExecMutex.Unlock()
	return e.ExecFatalErrorOccurred
}

func (e *BaseEnv) SetExecFatalErrorOccurred(err error) {
	e.ExecMutex.Lock()
	defer e.ExecMutex.Unlock()
	e.ExecFatalErrorOccurred = err
}

func (e *BaseEnv) GetCheckRunID() *int64 {
	return e.CheckRunID
}

func NewTypeEnv(e Env) TypeEnv {
	builtInsType := make(map[string]lang.Type)
	for builtInName, builtInFunction := range e.GetBuiltIns().Functions {
		builtInsType[builtInName] = builtInFunction.Type
	}

	for builtInName, builtInAction := range e.GetBuiltIns().Actions {
		builtInsType[builtInName] = builtInAction.Type
	}

	for valueName, value := range e.GetRegisterMap() {
		builtInsType[valueName] = value.Type()
	}

	return TypeEnv(builtInsType)
}

func NewEvalEnv(
	ctx context.Context,
	logger *logrus.Entry,
	dryRun bool,
	githubClient *gh.GithubClient,
	codeHostClient *codehost.CodeHostClient,
	collector collector.Collector,
	targetEntity *entities.TargetEntity,
	eventPayload interface{},
	builtIns *BuiltIns,
	checkRunID *int64,
) (Env, error) {
	registerMap := RegisterMap(make(map[string]lang.Value))
	report := &Report{Actions: make([]string, 0)}

	var wg sync.WaitGroup
	var mu sync.Mutex

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
		CodeHostClient:           codeHostClient,
		ExecWaitGroup:            &wg,
		ExecMutex:                &mu,
		CheckRunID:               checkRunID,
	}

	switch targetEntity.Kind {
	case entities.Issue:
		issue, _, err := githubClient.GetIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		if err != nil {
			return nil, err
		}

		input.Target = target.NewIssueTarget(ctx, targetEntity, githubClient, issue)
	case entities.PullRequest:
		pullRequest, err := codeHostClient.GetPullRequest(ctx, fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo), int64(targetEntity.Number))
		if err != nil {
			return nil, err
		}

		pullRequestTarget, err := target.NewPullRequestTarget(ctx, targetEntity, githubClient, codeHostClient, pullRequest)
		if err != nil {
			return nil, err
		}
		input.Target = pullRequestTarget
	}

	return input, nil
}
