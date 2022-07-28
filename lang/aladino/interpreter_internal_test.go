// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
)

func TestBuildGroupAST_WhenGroupTypeFilterIsSetAndParseFails(t *testing.T) {
	groupName := "senior-developers"

	gotExpr, err := buildGroupAST(
		engine.GroupTypeFilter,
		fmt.Sprintf("$group(\"%v\")", groupName),
		"dev",
		"$hasFileExtensions(",
	)

	assert.Nil(t, gotExpr)
	assert.EqualError(t, err, "parse error: failed to build AST on input $hasFileExtensions(")
}

func TestBuildGroupAST_WhenGroupTypeFilterIsSet(t *testing.T) {
	groupName := "senior-developers"

	gotExpr, err := buildGroupAST(
		engine.GroupTypeFilter,
		fmt.Sprintf("$group(\"%v\")", groupName),
		"dev",
		"$hasFileExtensions([\".ts\"])",
	)

	wantExpr := BuildFunctionCall(
		BuildVariable("filter"),
		[]Expr{
			BuildFunctionCall(
				BuildVariable("organization"),
				[]Expr{},
			),
			BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildFunctionCall(
					BuildVariable("hasFileExtensions"),
					[]Expr{
						BuildArray([]Expr{BuildStringConst(".ts")}),
					},
				),
			),
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestBuildGroupAST_WhenGroupTypeFilterIsNotSet(t *testing.T) {
	devName := "jane"

	gotExpr, err := buildGroupAST(
		engine.GroupTypeStatic,
		fmt.Sprintf("[\"%v\"]", devName),
		"",
		"",
	)

	wantExpr := BuildArray([]Expr{BuildStringConst(devName)})

	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestEvalGroup_WhenTypeInferenceFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, &BuiltIns{})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed %v", err))
	}

	expr, err := Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalGroup(mockedEnv, expr)

	assert.EqualError(t, err, "type inference failed")
}

func TestEvalGroup_WhenExpressionIsNotValidGroup(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, &BuiltIns{})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed %v", err))
	}

	expr, err := Parse("true")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalGroup(mockedEnv, expr)

	assert.EqualError(t, err, "expression is not a valid group")
}

func TestEvalGroup(t *testing.T) {
	devName := "jane"

	builtIns := &BuiltIns{
		Functions: map[string]*BuiltInFunction{
			"group": {
				Type: BuildFunctionType([]Type{BuildStringType()}, BuildArrayOfType(BuildStringType())),
				Code: func(e Env, args []Value) (Value, error) {
					return BuildArrayValue([]Value{BuildStringValue(devName)}), nil
				},
			},
		},
	}

	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, builtIns)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed %v", err))
	}

	expr, err := Parse("$group(\"\")")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	gotVal, err := evalGroup(mockedEnv, expr)

	wantVal := BuildArrayValue([]Value{BuildStringValue(devName)})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestProcessGroup_WhenBuildGroupASTFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, &BuiltIns{})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	errExpr := "$group("
	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		errExpr,
		"",
		"",
	)

	assert.EqualError(t, err, fmt.Sprintf("ProcessGroup:buildGroupAST: parse error: failed to build AST on input %v", errExpr))
}

func TestProcessGroup_WhenEvalGroupFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, &BuiltIns{})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	errExpr := "true"
	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		errExpr,
		"",
		"",
	)

	assert.EqualError(t, err, "ProcessGroup:evalGroup expression is not a valid group")
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSet(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, &BuiltIns{})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"
	devName := "jane"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		fmt.Sprintf("[\"%v\"]", devName),
		"",
		"",
	)

	gotVal := mockedEnv.GetRegisterMap()[groupName]

	wantVal := BuildArrayValue([]Value{
		BuildStringValue(devName),
	})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildInternalLabelID(t *testing.T) {
	labelID := "label_id"

	wantVal := fmt.Sprintf("@label:%v", labelID)

	gotVal := BuildInternalLabelID(labelID)

	assert.Equal(t, wantVal, gotVal)
}

func TestProcessLabel(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	labelID := "label_id"
	labelName := "label_name"
	err = mockedInterpreter.ProcessLabel(labelID, labelName)

	internalLabelID := fmt.Sprintf("@label:%v", labelID)
	gotVal := mockedEnv.GetRegisterMap()[internalLabelID]

	wantVal := BuildStringValue(labelName)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildInternalRuleName(t *testing.T) {
	ruleName := "rule_name"

	wantVal := fmt.Sprintf("@rule:%v", ruleName)

	gotVal := BuildInternalRuleName(ruleName)

	assert.Equal(t, wantVal, gotVal)
}

func TestProcessRule(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	ruleName := "rule_name"
	spec := "1 == 1"
	err = mockedInterpreter.ProcessRule(ruleName, spec)

	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	gotVal := mockedEnv.GetRegisterMap()[internalRuleName]

	wantVal := BuildStringValue(spec)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEvalExpr_WhenParseFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	gotVal, err := EvalExpr(mockedEnv, "", "1 ==")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "parse error: failed to build AST on input 1 ==")
}

func TestEvalExpr_WhenTypeInferenceFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	gotVal, err := EvalExpr(mockedEnv, "", "1 == \"a\"")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestEvalExpr_WhenExprIsNotBoolType(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	gotVal, err := EvalExpr(mockedEnv, "", "1")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "expression 1 is not a condition")
}

func TestEvalExpr(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	gotVal, err := EvalExpr(mockedEnv, "", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestEvalExpr_OnInterpreter(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	gotVal, err := mockedInterpreter.EvalExpr("", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestExecProgram_WhenExecStatementFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	program := &engine.Program{
		Statements: []*engine.Statement{
			{
				Code:     "$action()",
				Metadata: nil,
			},
		},
	}

	err = mockedInterpreter.ExecProgram(program)

	assert.EqualError(t, err, "no type for built-in action. Please check if the mode in the reviewpad.yml file supports it")
}

func TestExecProgram(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: BuildFunctionType([]Type{BuildStringType()}, nil),
				Code: func(e Env, args []Value) error {
					return nil
				},
			},
		},
	}

	mockedEnv, err := MockDefaultEnvWithBuiltIns(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{},
			),
			mock.WithRequestMatch(
				mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
				[]*github.Label{
					{Name: github.String("test")},
				},
			),
		},
		nil,
		builtIns,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statementWorkflowName := "test"
	statementRule := "testRule"
	statementCode := "$addLabel(\"test\")"
	statement := &engine.Statement{
		Code: statementCode,
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name: statementWorkflowName,
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: statementRule},
			},
		},
	}

	program := &engine.Program{
		Statements: []*engine.Statement{
			statement,
		},
	}

	err = mockedInterpreter.ExecProgram(program)

	gotVal := mockedEnv.GetReport().WorkflowDetails[statementWorkflowName]

	wantVal := ReportWorkflowDetails{
		Name: statementWorkflowName,
		Rules: map[string]bool{
			statementRule: true,
		},
		Actions: []string{statementCode},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestExecStatement_WhenParseFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := &engine.Statement{
		Code: "$addLabel(",
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name: "test",
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: "testRule"},
			},
		},
	}

	err = mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "parse error: failed to build AST on input $addLabel(")
}

func TestExecStatement_WhenTypeCheckExecFails(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: BuildFunctionType([]Type{BuildStringType()}, nil),
				Code: func(e Env, args []Value) error {
					return nil
				},
			},
		},
	}

	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, builtIns)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := &engine.Statement{
		Code: "$addLabel(1)",
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name: "test",
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: "testRule"},
			},
		},
	}

	err = mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "type inference failed: mismatch in arg types on addLabel")
}

func TestExecStatement_WhenActionExecFails(t *testing.T) {
	devName := "jane"
	builtIns := &BuiltIns{
		Functions: map[string]*BuiltInFunction{
			"author": {
				Type: BuildFunctionType([]Type{}, BuildArrayOfType(BuildStringType())),
				Code: func(e Env, args []Value) (Value, error) {
					return BuildArrayValue([]Value{BuildStringValue(devName)}), nil
				},
			},
		},
	}

	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, builtIns)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := &engine.Statement{
		Code: "$author()",
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name: "test",
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: "testRule"},
			},
		},
	}

	err = mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "exec: author not found. are you sure this is a built-in function?")
}

func TestExecStatement(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: BuildFunctionType([]Type{BuildStringType()}, nil),
				Code: func(e Env, args []Value) error {
					return nil
				},
			},
		},
	}

	mockedEnv, err := MockDefaultEnvWithBuiltIns(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{},
			),
			mock.WithRequestMatch(
				mock.PostReposIssuesLabelsByOwnerByRepoByIssueNumber,
				[]*github.Label{
					{Name: github.String("test")},
				},
			),
		},
		nil,
		builtIns,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnvWithBuiltIns failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statementWorkflowName := "test"
	statementRule := "testRule"
	statementCode := "$addLabel(\"test\")"
	statement := &engine.Statement{
		Code: statementCode,
		Metadata: &engine.Metadata{
			Workflow: engine.PadWorkflow{
				Name: statementWorkflowName,
			},
			TriggeredBy: []engine.PadWorkflowRule{
				{Rule: statementRule},
			},
		},
	}

	err = mockedInterpreter.ExecStatement(statement)

	gotVal := mockedEnv.GetReport().WorkflowDetails[statementWorkflowName]

	wantVal := ReportWorkflowDetails{
		Name: statementWorkflowName,
		Rules: map[string]bool{
			statementRule: true,
		},
		Actions: []string{statementCode},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestReport_WhenFindReportCommentFails(t *testing.T) {
	mockedPullRequest := GetDefaultMockPullRequestDetailsWith(&github.PullRequest{
		User: &github.User{Login: github.String("john")},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String("john"),
				},
				Name: github.String("default-mock-repo"),
			},
			Ref: github.String("master"),
		},
	})
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.EqualError(t, err, "[report] error getting issues mock response not found for /repos/john/default-mock-repo/issues/6/comments")
}

func TestReport_OnSilentMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var isDeletedCommentRequested bool
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// When a pull request has a reviewpad comment then this comment is deleted
					isDeletedCommentRequested = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.True(t, isDeletedCommentRequested)
}

func TestReport_OnSilentMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	var isDeletedCommentRequested bool
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.DeleteReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// When a pull request has a reviewpad comment then this comment is deleted
					isDeletedCommentRequested = true
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.False(t, isDeletedCommentRequested)
}

func TestReport_OnVerboseMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	var addedComment string
	commentToBeAdded := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{},
			),
			mock.WithRequestMatchHandler(
				mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := github.IssueComment{}

					json.Unmarshal(rawBody, &body)

					addedComment = *body.Body
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.VERBOSE_MODE, false)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeAdded, addedComment)
}

func TestReport_OnVerboseMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var updatedComment string
	commentUpdated := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated"
	mockedEnv, err := MockDefaultEnv(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String("<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n**:information_source: Messages**\n* Changes the README.md"),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := ioutil.ReadAll(r.Body)
					body := github.IssueComment{}

					json.Unmarshal(rawBody, &body)

					updatedComment = *body.Body
				}),
			),
		},
		nil,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err = mockedInterpreter.Report(engine.VERBOSE_MODE, false)

	assert.Nil(t, err)
	assert.Equal(t, commentUpdated, updatedComment)
}

func TestNewInterpreter_WhenNewEvalEnvFails(t *testing.T) {
	ctx := context.Background()
	failMessage := "GetPullRequestFilesRequestFail"
	client := github.NewClient(
		mock.NewMockedHTTPClient(
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		),
	)
	// TODO: Ideally, we should not have nil arguments in the call to NewInterpreter
	gotInterpreter, err := NewInterpreter(
		ctx,
		false,
		client,
		nil,
		nil,
		GetDefaultMockPullRequestDetails(),
		nil,
		nil,
	)

	assert.Nil(t, gotInterpreter)
	assert.NotNil(t, err)
}

func TestNewInterpreter(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("MockDefaultEnv failed: %v", err))
	}

	wantInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	gotInterpreter, err := NewInterpreter(
		mockedEnv.GetCtx(),
		mockedEnv.GetDryRun(),
		mockedEnv.GetClient(),
		mockedEnv.GetClientGQL(),
		mockedEnv.GetCollector(),
		mockedEnv.GetPullRequest(),
		mockedEnv.GetEventPayload(),
		mockedEnv.GetBuiltIns(),
	)

	assert.Nil(t, err)
	assert.Equal(t, wantInterpreter, gotInterpreter)
}
