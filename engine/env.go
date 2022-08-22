// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/lang/aladino/target"
)

type GroupKind string
type GroupType string

type ExitStatus int

const GroupKindDeveloper GroupKind = "developer"
const GroupTypeStatic GroupType = "static"
const GroupTypeFilter GroupType = "filter"

const ExitStatusSuccess ExitStatus = 0
const ExitStatusFailure ExitStatus = 1

type Interpreter interface {
	ProcessGroup(name string, kind GroupKind, typeOf GroupType, expr, paramExpr, whereExpr string) error
	ProcessLabel(id, name string) error
	ProcessRule(name, spec string) error
	EvalExpr(kind, expr string) (bool, error)
	ExecProgram(program *Program) (ExitStatus, error)
	ExecStatement(statement *Statement) error
	Report(mode string, safeMode bool) error
}

type Env struct {
	Ctx          context.Context
	DryRun       bool
	GithubClient *gh.GithubClient
	Collector    collector.Collector
	PullRequest  *github.PullRequest
	EventPayload interface{}
	Interpreter  Interpreter
	TargetEntity *handler.TargetEntity
	Issue        *github.Issue
	Target       target.Target
}

func NewEvalEnv(
	ctx context.Context,
	dryRun bool,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	targetEntity *handler.TargetEntity,
	eventPayload interface{},
	interpreter Interpreter,
) (*Env, error) {
	input := &Env{
		Ctx:          ctx,
		DryRun:       dryRun,
		GithubClient: githubClient,
		Collector:    collector,
		EventPayload: eventPayload,
		Interpreter:  interpreter,
		TargetEntity: targetEntity,
	}

	if targetEntity.Kind == handler.PullRequest {
		pullRequest, _, err := githubClient.GetPullRequest(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
		if err != nil {
			return nil, err
		}

		input.PullRequest = pullRequest
		input.Target = target.NewPullRequestTarget(ctx, targetEntity, githubClient, pullRequest)

		return input, nil
	}

	issue, _, err := githubClient.GetIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
	if err != nil {
		return nil, err
	}

	input.Issue = issue
	input.Target = target.NewIssueTarget(ctx, targetEntity, githubClient, issue)

	return input, nil
}
