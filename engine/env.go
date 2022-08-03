// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/shurcooL/githubv4"
)

type GroupKind string
type GroupType string

const GroupKindDeveloper GroupKind = "developer"
const GroupTypeStatic GroupType = "static"
const GroupTypeFilter GroupType = "filter"

type Interpreter interface {
	ProcessGroup(name string, kind GroupKind, typeOf GroupType, expr, paramExpr, whereExpr string) error
	ProcessLabel(id, name string) error
	ProcessRule(name, spec string) error
	EvalExpr(kind, expr string) (bool, error)
	ExecProgram(program *Program) error
	ExecStatement(statement *Statement) error
	Report(mode string, safeMode bool) error
	CheckForFatal()
}

type Env struct {
	Ctx          context.Context
	DryRun       bool
	Client       *github.Client
	ClientGQL    *githubv4.Client
	Collector    collector.Collector
	PullRequest  *github.PullRequest
	EventPayload interface{}
	Interpreter  Interpreter
}

func NewEvalEnv(
	ctx context.Context,
	dryRun bool,
	client *github.Client,
	clientGQL *githubv4.Client,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	eventPayload interface{},
	interpreter Interpreter,
) (*Env, error) {
	input := &Env{
		Ctx:          ctx,
		DryRun:       dryRun,
		Client:       client,
		ClientGQL:    clientGQL,
		Collector:    collector,
		PullRequest:  pullRequest,
		EventPayload: eventPayload,
		Interpreter:  interpreter,
	}

	return input, nil
}
