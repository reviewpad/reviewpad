// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"

	"github.com/reviewpad/go-lib/entities"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/sirupsen/logrus"
)

type GroupKind string
type GroupType string
type ExitStatus int

const GroupKindDeveloper GroupKind = "developer"
const GroupTypeStatic GroupType = "static"
const GroupTypeFilter GroupType = "filter"

const ExitStatusSuccess ExitStatus = 0
const ExitStatusFailure ExitStatus = 1

type CheckState string

const (
	CheckStateError   CheckState = "error"
	CheckStateFailure CheckState = "failure"
	CheckStatePending CheckState = "pending"
	CheckStateSuccess CheckState = "success"
)

type Interpreter interface {
	ProcessGroup(name string, kind GroupKind, typeOf GroupType, expr, paramExpr, whereExpr string) error
	ProcessLabel(id, name string) error
	ProcessRule(name, spec string) error
	EvalExpr(kind, expr string) (bool, error)
	ExecProgram(program *Program) (ExitStatus, error)
	ExecStatement(statement *Statement) error
	Report(mode string, safeMode bool) error
	ReportMetrics() error
	GetCheckRunConclusion() string
	ProcessList(name, expr string) error
}

type Env struct {
	Ctx          context.Context
	DryRun       bool
	GithubClient *gh.GithubClient
	Collector    collector.Collector
	Interpreter  Interpreter
	TargetEntity *entities.TargetEntity
	EventDetails *entities.EventDetails
	Logger       *logrus.Entry
}

func NewEvalEnv(
	ctx context.Context,
	logger *logrus.Entry,
	dryRun bool,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	targetEntity *entities.TargetEntity,
	interpreter Interpreter,
	eventDetails *entities.EventDetails,
) (*Env, error) {
	input := &Env{
		Ctx:          ctx,
		DryRun:       dryRun,
		GithubClient: githubClient,
		Collector:    collector,
		Interpreter:  interpreter,
		TargetEntity: targetEntity,
		EventDetails: eventDetails,
		Logger:       logger,
	}

	return input, nil
}
