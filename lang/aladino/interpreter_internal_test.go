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

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v45/github"
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, DefaultMockEventPayload, controller)

	expr, err := Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalGroup(mockedEnv, expr)

	assert.EqualError(t, err, "type inference failed")
}

func TestEvalGroup_WhenExpressionIsNotValidGroup(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, DefaultMockEventPayload, controller)

	expr, err := Parse("true")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalGroup(mockedEnv, expr)

	assert.EqualError(t, err, "expression is not a valid group")
}

func TestEvalGroup(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, DefaultMockEventPayload, controller)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	errExpr := "$group("
	err := mockedInterpreter.ProcessGroup(
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	errExpr := "true"
	err := mockedInterpreter.ProcessGroup(
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"
	devName := "jane"

	err := mockedInterpreter.ProcessGroup(
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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	labelID := "label_id"
	labelName := "label_name"
	err := mockedInterpreter.ProcessLabel(labelID, labelName)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	ruleName := "rule_name"
	spec := "1 == 1"
	err := mockedInterpreter.ProcessRule(ruleName, spec)

	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	gotVal := mockedEnv.GetRegisterMap()[internalRuleName]

	wantVal := BuildStringValue(spec)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEvalExpr_WhenParseFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	gotVal, err := EvalExpr(mockedEnv, "", "1 ==")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "parse error: failed to build AST on input 1 ==")
}

func TestEvalExpr_WhenTypeInferenceFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	gotVal, err := EvalExpr(mockedEnv, "", "1 == \"a\"")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestEvalExpr_WhenExprIsNotBoolType(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	gotVal, err := EvalExpr(mockedEnv, "", "1")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "expression 1 is not a condition")
}

func TestEvalExpr(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	gotVal, err := EvalExpr(mockedEnv, "", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestEvalExpr_OnInterpreter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	gotVal, err := mockedInterpreter.EvalExpr("", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestExecProgram_WhenExecStatementFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$action()")
	statements := []*engine.Statement{statement}
	program := engine.BuildProgram(statements)

	exitStatus, err := mockedInterpreter.ExecProgram(program)

	assert.Equal(t, engine.ExitStatusFailure, exitStatus)
	assert.EqualError(t, err, "no type for built-in action. Please check if the mode in the reviewpad.yml file supports it")
}

func TestExecProgram(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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

	mockedEnv := MockDefaultEnv(
		t,
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
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statCode := "$addLabel(\"test\")"

	statement := engine.BuildStatement(statCode)
	program := engine.BuildProgram([]*engine.Statement{statement})

	exitStatus, err := mockedInterpreter.ExecProgram(program)

	gotVal := mockedEnv.GetReport()

	wantVal := &Report{
		Actions: []string{statCode},
	}

	assert.Nil(t, err)
	assert.Equal(t, engine.ExitStatusSuccess, exitStatus)
	assert.Equal(t, wantVal, gotVal)
}

func TestExecStatement_WhenParseFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$addLabel(")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "parse error: failed to build AST on input $addLabel(")
}

func TestExecStatement_WhenTypeCheckExecFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$addLabel(1)")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "type inference failed: mismatch in arg types on addLabel")
}

func TestExecStatement_WhenActionExecFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, DefaultMockEventPayload, controller)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$author()")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "exec: author not found. are you sure this is a built-in function?")
}

func TestExecStatement(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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

	mockedEnv := MockDefaultEnv(
		t,
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
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statCode := "$addLabel(\"test\")"

	statement := engine.BuildStatement(statCode)

	err := mockedInterpreter.ExecStatement(statement)

	gotVal := mockedEnv.GetReport()

	wantVal := &Report{
		Actions: []string{statCode},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestReport_WhenFindReportCommentFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

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
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Write(mock.MustMarshal(mockedPullRequest))
				}),
			),
		},
		nil,
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.EqualError(t, err, "[report] error getting issues mock response not found for /repos/john/default-mock-repo/issues/6/comments")
}

func TestReport_OnSilentMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var isDeletedCommentRequested bool
	mockedEnv := MockDefaultEnv(
		t,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.True(t, isDeletedCommentRequested)
}

func TestReport_OnSilentMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var isDeletedCommentRequested bool
	mockedEnv := MockDefaultEnv(
		t,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.False(t, isDeletedCommentRequested)
}

func TestReport_OnVerboseMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var addedComment string
	commentToBeAdded := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Executed actions**\n```yaml\n```\n"
	mockedEnv := MockDefaultEnv(
		t,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.VERBOSE_MODE, false)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeAdded, addedComment)
}

func TestReport_OnVerboseMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var updatedComment string
	commentUpdated := "<!--@annotation-reviewpad-report-->\n**Reviewpad Report**\n\n:scroll: **Executed actions**\n```yaml\n```\n"
	mockedEnv := MockDefaultEnv(
		t,
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
		MockBuiltIns(),
		DefaultMockEventPayload,
		controller,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.VERBOSE_MODE, false)

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
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), DefaultMockEventPayload, controller)

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
