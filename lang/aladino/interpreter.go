// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
	"github.com/shurcooL/githubv4"
)

type Interpreter struct {
	Env Env
}

func execLog(val string) {
	log.Println(fmtio.Sprint("aladino", val))
}

func execLogf(format string, a ...interface{}) {
	log.Println(fmtio.Sprintf("aladino", format, a...))
}

func buildGroupAST(typeOf engine.GroupType, expr, paramExpr, whereExpr string) (Expr, error) {
	if typeOf == engine.GroupTypeFilter {
		whereExprAST, err := Parse(whereExpr)
		if err != nil {
			return nil, err
		}

		return BuildFilter(paramExpr, whereExprAST)
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

func (i *Interpreter) ProcessGroup(groupName string, kind engine.GroupKind, typeOf engine.GroupType, expr, paramExpr, whereExpr string) error {
	exprAST, err := buildGroupAST(typeOf, expr, paramExpr, whereExpr)
	if err != nil {
		return fmt.Errorf("ProcessGroup:buildGroupAST: %v", err)
	}

	value, err := evalGroup(i.Env, exprAST)
	if err != nil {
		return fmt.Errorf("ProcessGroup:evalGroup %v", err)
	}

	i.Env.GetRegisterMap()[groupName] = value
	return nil
}

func BuildInternalLabelID(id string) string {
	return fmt.Sprintf("@label:%v", id)
}

func (i *Interpreter) ProcessLabel(id, name string) error {
	internalLabelID := BuildInternalLabelID(id)

	i.Env.GetRegisterMap()[internalLabelID] = BuildStringValue(name)
	return nil
}

func BuildInternalRuleName(name string) string {
	return fmt.Sprintf("@rule:%v", name)
}

func (i *Interpreter) ProcessRule(name, spec string) error {
	internalRuleName := BuildInternalRuleName(name)

	i.Env.GetRegisterMap()[internalRuleName] = BuildStringValue(spec)
	return nil
}

func EvalExpr(env Env, kind, expr string) (bool, error) {
	exprAST, err := Parse(expr)
	if err != nil {
		return false, err
	}

	exprType, err := TypeInference(env, exprAST)
	if err != nil {
		return false, err
	}

	if exprType.Kind() != BOOL_TYPE {
		return false, fmt.Errorf("expression %v is not a condition", expr)
	}

	return EvalCondition(env, exprAST)
}

func (i *Interpreter) EvalExpr(kind, expr string) (bool, error) {
	return EvalExpr(i.Env, kind, expr)
}

func (i *Interpreter) ExecProgram(program *engine.Program) error {
	execLog("executing program")

	for _, statement := range program.Statements {
		err := i.ExecStatement(statement)
		if err != nil {
			return err
		}
	}

	execLog("execution done")

	return nil
}

func (i *Interpreter) ExecStatement(statement *engine.Statement) error {
	statRaw := statement.Code
	statAST, err := Parse(statRaw)
	if err != nil {
		return err
	}

	execStatAST, err := TypeCheckExec(i.Env, statAST)
	if err != nil {
		return err
	}

	m := regexp.MustCompile("\\$fail\\(\"(.*)\"\\)")

	isFailAction := m.FindString(statement.Code) != ""

	if !i.Env.GetDryRun() {
		// The report is generated after the program is executed.
		// $fail built-in action will fail the GitHub action before the report is generated.
		// In this case, when we encounter the $fail built-in action we need to first generate the report before executing the action.
		if isFailAction {
			i.Env.GetReport().addToReport(statement)

			err := i.Report()
			if err != nil {
				return err
			}
		}

		err = execStatAST.exec(i.Env)
		if err != nil {
			return err
		}
	}

	i.Env.GetReport().addToReport(statement)

	execLogf("\taction %v executed", statRaw)
	return nil
}

func (i *Interpreter) Report() error {
	execLog("generating report")

	env := i.Env

	useSafeModeHeader := env.GetReport().Settings.UseSafeModeHeader
	createReportComment := env.GetReport().Settings.CreateReportComment

	var err error

	comment, err := FindReportComment(env)
	if err != nil {
		return err
	}

	if !createReportComment {
		if comment != nil {
			return DeleteReportComment(env, *comment.ID)
		}
		return nil
	}

	report := buildReport(useSafeModeHeader, env.GetReport())

	if comment == nil {
		return AddReportComment(env, report)
	}

	return UpdateReportComment(env, *comment.ID, report)

}

func NewInterpreter(
	ctx context.Context,
	dryRun bool,
	safeMode bool,
	createReportComment bool,
	gitHubClient *github.Client,
	gitHubClientGQL *githubv4.Client,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	eventPayload interface{},
	builtIns *BuiltIns,

) (engine.Interpreter, error) {
	evalEnv, err := NewEvalEnv(ctx, dryRun, safeMode, createReportComment, gitHubClient, gitHubClientGQL, collector, pullRequest, eventPayload, builtIns)
	if err != nil {
		return nil, err
	}

	return &Interpreter{
		Env: evalEnv,
	}, nil
}
