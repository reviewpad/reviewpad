// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"fmt"
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/reviewpad/reviewpad/v4/engine"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/sirupsen/logrus"
)

type Interpreter struct {
	Env Env
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

func evalGroup(env Env, expr Expr) (lang.Value, error) {
	exprType, err := TypeInference(env, expr)
	if err != nil {
		return nil, err
	}

	if exprType.Kind() != lang.ARRAY_TYPE && exprType.Kind() != lang.ARRAY_OF_TYPE {
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

func (i *Interpreter) ProcessIterable(expr string) (lang.Value, error) {
	exprAST, err := Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("ProcessIterable:Parse: %v", err)
	}

	exprType, err := TypeInference(i.Env, exprAST)
	if err != nil {
		return nil, err
	}

	if exprType.Kind() != lang.ARRAY_TYPE && exprType.Kind() != lang.ARRAY_OF_TYPE && exprType.Kind() != lang.DICTIONARY_TYPE {
		return nil, fmt.Errorf("expression is not a valid iterable")
	}

	return Eval(i.Env, exprAST)
}

func (i *Interpreter) ProcessDictionary(name string, dictionary map[string]string) error {
	processedDictionary := make(map[string]lang.Value)

	for key, expr := range dictionary {
		exprAST, err := Parse(expr)
		if err != nil {
			return fmt.Errorf("ProcessDictionary:Parse: %v", err)
		}

		value, err := Eval(i.Env, exprAST)
		if err != nil {
			return fmt.Errorf("ProcessDictionary:Eval: %v", err)
		}

		processedDictionary[key] = value
	}

	i.Env.GetRegisterMap()[fmt.Sprintf("@dictionary:%s", name)] = lang.BuildDictionaryValue(processedDictionary)

	return nil
}

func (i *Interpreter) StoreTemporaryVariable(name string, value lang.Value) {
	i.Env.GetRegisterMap()[BuildInternalTemporaryVariableName(name)] = value
}

func BuildInternalTemporaryVariableName(name string) string {
	return fmt.Sprintf("@variable:%s", name)
}

func BuildInternalLabelID(id string) string {
	return fmt.Sprintf("@label:%v", id)
}

func (i *Interpreter) ProcessLabel(id, name string) error {
	internalLabelID := BuildInternalLabelID(id)

	i.Env.GetRegisterMap()[internalLabelID] = lang.BuildStringValue(name)
	return nil
}

func BuildInternalRuleName(name string) string {
	return fmt.Sprintf("@rule:%v", name)
}

func (i *Interpreter) ProcessRule(name, spec string) error {
	internalRuleName := BuildInternalRuleName(name)

	i.Env.GetRegisterMap()[internalRuleName] = lang.BuildStringValue(spec)
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

	if exprType.Kind() != lang.BOOL_TYPE {
		return false, fmt.Errorf("expression %v is not a condition", expr)
	}

	return EvalCondition(env, exprAST)
}

func (i *Interpreter) EvalExpr(kind, expr string) (bool, error) {
	return EvalExpr(i.Env, kind, expr)
}

func (i *Interpreter) ExecProgram(program *engine.Program) (engine.ExitStatus, error) {
	retStatus := engine.ExitStatusSuccess
	var retErr error

	for _, statement := range program.GetProgramStatements() {
		err := i.ExecStatement(statement)
		if err != nil {
			retStatus = engine.ExitStatusFailure
			retErr = err
			break
		}

		hasFatalError := len(i.Env.GetBuiltInsReportedMessages()[SEVERITY_FATAL]) > 0
		if hasFatalError {
			i.Env.GetLogger().Info("execution stopped")
			retStatus = engine.ExitStatusFailure
			retErr = nil
			break
		}
	}

	i.Env.GetExecWaitGroup().Wait()

	if retStatus == engine.ExitStatusFailure {
		return retStatus, retErr
	}

	if i.Env.GetExecFatalErrorOccurred() != nil {
		return engine.ExitStatusFailure, i.Env.GetExecFatalErrorOccurred()
	}

	return retStatus, retErr
}

func (i *Interpreter) ExecStatement(statement *engine.Statement) error {
	statRaw := statement.GetStatementCode()
	statAST, err := Parse(statRaw)
	if err != nil {
		return err
	}

	execStatAST, err := TypeCheckExec(i.Env, statAST)
	if err != nil {
		return err
	}

	if !i.Env.GetDryRun() {
		err = execStatAST.exec(i.Env)
		if err != nil {
			return err
		}

		// update the target
		targetEntity := i.Env.GetTarget().GetTargetEntity()
		githubClient := i.Env.GetGithubClient()
		codeHostClient := i.Env.GetCodeHostClient()
		ctx := i.Env.GetCtx()
		switch targetEntity.Kind {
		case entities.Issue:
			issue, _, err := githubClient.GetIssue(ctx, targetEntity.Owner, targetEntity.Repo, targetEntity.Number)
			if err != nil {
				return err
			}

			i.Env.(*BaseEnv).Target = target.NewIssueTarget(ctx, targetEntity, githubClient, issue)
		case entities.PullRequest:
			pullRequest, err := codeHostClient.GetPullRequest(ctx, fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo), int64(targetEntity.Number))
			if err != nil {
				return err
			}

			pullRequestTarget, err := target.NewPullRequestTarget(ctx, targetEntity, githubClient, codeHostClient, pullRequest)
			if err != nil {
				return err
			}
			i.Env.(*BaseEnv).Target = pullRequestTarget
		}
	}

	i.Env.GetReport().addToReport(statement)

	i.Env.GetLogger().Infof("action `%v` executed", statRaw)
	return nil
}

func (i *Interpreter) Report(mode string, safeMode bool) error {
	i.Env.GetLogger().Info("generating report")

	if mode == "" {
		// By default mode is silent
		mode = engine.SILENT_MODE
	}

	env := i.Env

	var err error

	comment, err := FindReportCommentByAnnotation(env, ReviewpadReportCommentAnnotation)
	if err != nil {
		return err
	}

	reportComments := env.GetBuiltInsReportedMessages()

	// Since fail messages aren't supposed to be reported, we remove them from the report
	delete(reportComments, SEVERITY_FAIL)

	if mode == engine.SILENT_MODE && len(reportComments) == 0 && !safeMode {
		if comment != nil {
			return DeleteReportComment(env, *comment.ID)
		}
		return nil
	}

	report := buildReport(mode, safeMode, reportComments, env.GetReport())

	if comment == nil {
		return AddReportComment(env, report)
	}

	return UpdateReportComment(env, *comment.ID, report)
}

func (i *Interpreter) ReportMetrics() error {
	targetEntity := i.Env.GetTarget().GetTargetEntity()
	owner := targetEntity.Owner
	prNum := targetEntity.Number
	repo := targetEntity.Repo
	ctx := i.Env.GetCtx()
	pr := i.Env.GetTarget().(*target.PullRequestTarget).PullRequest

	if !pr.IsMerged {
		return nil
	}

	report := strings.Builder{}

	firstCommitDate, firstReviewDate, err := i.Env.GetGithubClient().GetFirstCommitAndReviewDate(ctx, owner, repo, prNum)
	if err != nil {
		return err
	}

	if firstCommitDate != nil {
		report.WriteString(fmt.Sprintf("**💻 Coding Time**: %s", utils.ReadableTimeDiff(*firstCommitDate, pr.CreatedAt.AsTime())))
	}

	if firstReviewDate != nil && firstReviewDate.Before(pr.MergedAt.AsTime()) {
		report.WriteString(fmt.Sprintf("\n**🛻 Pickup Time**: %s", utils.ReadableTimeDiff(pr.CreatedAt.AsTime(), *firstReviewDate)))

		report.WriteString(fmt.Sprintf("\n**👀 Review Time**: %s", utils.ReadableTimeDiff(*firstReviewDate, pr.MergedAt.AsTime())))
	}

	if report.Len() > 0 {
		comment, err := FindReportCommentByAnnotation(i.Env, ReviewpadMetricReportCommentAnnotation)
		if err != nil {
			return err
		}

		r := fmt.Sprintf("%s\n## 📈 Pull Request Metrics\n%s", ReviewpadMetricReportCommentAnnotation, report.String())

		if comment == nil {
			err = AddReportComment(i.Env, r)
			if err != nil {
				return err
			}
			return nil
		}

		err = UpdateReportComment(i.Env, *comment.ID, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) GetCheckRunConclusion() string {
	return i.Env.GetCheckRunConclusion()
}

func NewInterpreter(
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
) (engine.Interpreter, error) {
	evalEnv, err := NewEvalEnv(ctx, logger, dryRun, githubClient, codeHostClient, collector, targetEntity, eventPayload, builtIns, checkRunID)
	if err != nil {
		return nil, err
	}

	return &Interpreter{
		Env: evalEnv,
	}, nil
}
