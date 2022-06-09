// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/collector"
	"github.com/reviewpad/reviewpad/engine"
	"github.com/reviewpad/reviewpad/utils/fmtio"
	"github.com/shurcooL/githubv4"
)

type Interpreter struct {
	Env Env
}

func (i *Interpreter) ProcessGroup(groupName string, kind engine.GroupKind, typeOf engine.GroupType, expr, paramExpr, whereExpr string) error {
	exprAST, _ := buildGroupAST(typeOf, expr, paramExpr, whereExpr)
	value, err := evalGroup(i.Env, exprAST)

	(*i.Env.GetRegisterMap())[groupName] = value

	return err
}

func buildGroupAST(typeOf engine.GroupType, expr, paramExpr, whereExpr string) (Expr, error) {
	if typeOf == engine.GroupTypeFilter {
		whereExprAST, err := Parse(whereExpr)
		if err != nil {
			return nil, err
		}

		return buildFilter(paramExpr, whereExprAST)
	} else {
		return Parse(expr)
	}
}

func evalGroup(env Env, expr Expr) (Value, error) {
	exprType, err := TypeInference(env, expr)
	if err != nil {
		return nil, err
	}

	if exprType.Kind() != ARRAY_TYPE && exprType.Kind() != ARRAY_OF_TYPE {
		return nil, fmt.Errorf("expression is not a valid group")
	}

	return Eval(env, expr)
}

func (i *Interpreter) EvalExpr(kind, expr string) (bool, error) {
	exprAST, err := Parse(expr)
	if err != nil {
		return false, err
	}

	exprType, err := TypeInference(i.Env, exprAST)
	if err != nil {
		return false, err
	}

	if exprType.Kind() != "BoolType" {
		return false, fmt.Errorf("expression %v is not a condition", expr)
	}

	cleanTemporaryVariables(i.Env)

	return EvalCondition(i.Env, exprAST)
}

func cleanTemporaryVariables(env Env) {
	registerMap := env.GetRegisterMap()
	for varName := range *registerMap {
		if varName[0] == '@' {
			delete(*registerMap, varName)
		}
	}
}

func execLog(val string) {
	log.Println(fmtio.Sprint("aladino", val))
}

func execLogf(format string, a ...interface{}) {
	log.Println(fmtio.Sprintf("aladino", format, a...))
}

func (i *Interpreter) ExecActions(program *[]string) error {
	execLog("executing actions:")

	for _, statRaw := range *program {
		statAST, err := Parse(statRaw)
		if err != nil {
			return err
		}

		execStatAST, err := TypeCheckExec(i.Env, statAST)
		if err != nil {
			return err
		}

		err = ExecAction(i.Env, execStatAST)
		if err != nil {
			return err
		}

		execLogf("\taction %v executed", statRaw)
	}

	execLog("execution done")

	return nil
}

func NewInterpreter(
	ctx context.Context,
	gitHubClient *github.Client,
	gitHubClientGQL *githubv4.Client,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	builtIns *BuiltIns,
) (engine.Interpreter, error) {
	evalEnv, err := NewEvalEnv(ctx, gitHubClient, gitHubClientGQL, collector, pullRequest, builtIns)
	if err != nil {
		return nil, err
	}

	return &Interpreter{
		Env: evalEnv,
	}, nil
}
