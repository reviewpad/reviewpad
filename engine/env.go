// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"

	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
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
	ReportMetrics() error
}

type Env struct {
	Ctx          context.Context
	DryRun       bool
	GithubClient *gh.GithubClient
	Collector    collector.Collector
	Interpreter  Interpreter
	TargetEntity *handler.TargetEntity
	EventData    *handler.EventData
}

func NewEvalEnv(
	ctx context.Context,
	dryRun bool,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	targetEntity *handler.TargetEntity,
	interpreter Interpreter,
	eventData *handler.EventData,
) (*Env, error) {
	input := &Env{
		Ctx:          ctx,
		DryRun:       dryRun,
		GithubClient: githubClient,
		Collector:    collector,
		Interpreter:  interpreter,
		TargetEntity: targetEntity,
		EventData:    eventData,
	}

	return input, nil
}
