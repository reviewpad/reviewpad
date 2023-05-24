// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/engine"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/utils"
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
				[]Expr{BuildTypedExpr(BuildVariable("dev"), lang.BuildStringType())},
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
	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, nil)

	expr, err := Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalList(mockedEnv, expr)

	assert.EqualError(t, err, "type inference failed")
}

func TestEvalGroup_WhenExpressionIsNotValidGroup(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, nil)

	expr, err := Parse("true")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	_, err = evalList(mockedEnv, expr)

	assert.EqualError(t, err, "expression is not a valid group")
}

func TestEvalGroup(t *testing.T) {
	devName := "jane"

	builtIns := &BuiltIns{
		Functions: map[string]*BuiltInFunction{
			"group": {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
				Code: func(e Env, args []lang.Value) (lang.Value, error) {
					return lang.BuildArrayValue([]lang.Value{lang.BuildStringValue(devName)}), nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
			},
		},
	}

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, nil)

	expr, err := Parse("$group(\"\")")
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("parse failed %v", err))
	}

	gotVal, err := evalList(mockedEnv, expr)

	wantVal := lang.BuildArrayValue([]lang.Value{lang.BuildStringValue(devName)})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestProcessGroup_WhenBuildGroupASTFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, nil)

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
	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, nil)

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
	mockedEnv := MockDefaultEnv(t, nil, nil, &BuiltIns{}, nil)

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

	wantVal := lang.BuildArrayValue([]lang.Value{
		lang.BuildStringValue(devName),
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
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	labelID := "label_id"
	labelName := "label_name"
	err := mockedInterpreter.ProcessLabel(labelID, labelName)

	internalLabelID := fmt.Sprintf("@label:%v", labelID)
	gotVal := mockedEnv.GetRegisterMap()[internalLabelID]

	wantVal := lang.BuildStringValue(labelName)

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
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	ruleName := "rule_name"
	spec := "1 == 1"
	err := mockedInterpreter.ProcessRule(ruleName, spec)

	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	gotVal := mockedEnv.GetRegisterMap()[internalRuleName]

	wantVal := lang.BuildStringValue(spec)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEvalExpr_WhenParseFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	gotVal, err := EvalExpr(mockedEnv, "", "1 ==")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "parse error: failed to build AST on input 1 ==")
}

func TestEvalExpr_WhenTypeInferenceFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	gotVal, err := EvalExpr(mockedEnv, "", "1 == \"a\"")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestEvalExpr_WhenExprIsNotBoolType(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	gotVal, err := EvalExpr(mockedEnv, "", "1")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "expression 1 is not a condition")
}

func TestEvalExpr(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	gotVal, err := EvalExpr(mockedEnv, "", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestEvalExpr_OnInterpreter(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	gotVal, err := mockedInterpreter.EvalExpr("", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestExecProgram_WhenExecStatementFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$action()")
	statements := []*engine.Statement{statement}
	program := engine.BuildProgram(statements)

	exitStatus, err := mockedInterpreter.ExecProgram(program)

	assert.Equal(t, engine.ExitStatusFailure, exitStatus)
	assert.EqualError(t, err, "no type for built-in action")
}

func TestExecProgram(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
				Code: func(e Env, args []lang.Value) error {
					return nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
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
			mock.WithRequestMatchHandler(mock.EndpointPattern{
				Pattern: "/repos/foobar/default-mock-repo/statuses/",
				Method:  "POST",
			}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{}`)
			})),
		},
		nil,
		builtIns,
		nil,
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

func TestExecAsyncProgram(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"robin": {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
				Code: func(e Env, args []lang.Value) error {
					return fmt.Errorf("robin error")
				},
				SupportedKinds:    []entities.TargetEntityKind{entities.PullRequest},
				RunAsynchronously: true,
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
			mock.WithRequestMatchHandler(mock.EndpointPattern{
				Pattern: "/repos/foobar/default-mock-repo/statuses/",
				Method:  "POST",
			}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				utils.MustWrite(w, `{}`)
			})),
		},
		nil,
		builtIns,
		nil,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statCode := "$robin(\"test\")"

	statement := engine.BuildStatement(statCode)
	program := engine.BuildProgram([]*engine.Statement{statement})

	exitStatus, err := mockedInterpreter.ExecProgram(program)

	gotVal := mockedEnv.GetReport()

	wantVal := &Report{
		Actions: []string{statCode},
	}

	assert.EqualError(t, err, "robin error")
	assert.Equal(t, engine.ExitStatusFailure, exitStatus)
	assert.Equal(t, wantVal, gotVal)
}

func TestExecStatement_WhenParseFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$addLabel(")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "parse error: failed to build AST on input $addLabel(")
}

func TestExecStatement_WhenTypeCheckExecFails(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
				Code: func(e Env, args []lang.Value) error {
					return nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
			},
		},
	}

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$addLabel(1)")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "type inference failed: mismatch in arg types on addLabel")
}

func TestExecStatement_WhenActionExecFails(t *testing.T) {
	devName := "jane"
	builtIns := &BuiltIns{
		Functions: map[string]*BuiltInFunction{
			"author": {
				Type: lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
				Code: func(e Env, args []lang.Value) (lang.Value, error) {
					return lang.BuildArrayValue([]lang.Value{lang.BuildStringValue(devName)}), nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
			},
		},
	}

	mockedEnv := MockDefaultEnv(t, nil, nil, builtIns, nil)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	statement := engine.BuildStatement("$author()")

	err := mockedInterpreter.ExecStatement(statement)

	assert.EqualError(t, err, "exec: author not found. are you sure this is a built-in function?")
}

func TestExecStatement(t *testing.T) {
	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			"addLabel": {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
				Code: func(e Env, args []lang.Value) error {
					return nil
				},
				SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
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
		nil,
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
	mockedCodeReview := GetDefaultMockPullRequestDetailsWith(&pbc.PullRequest{
		Author: &pbc.User{Login: "foobar"},
		Base: &pbc.Branch{
			Repo: &pbc.Repository{
				Owner: "foobar",
				Name:  "default-mock-repo",
			},
			Name: "master",
		},
	})
	mockedEnv := MockDefaultEnvWithPullRequestAndFiles(
		t,
		nil,
		nil,
		mockedCodeReview,
		GetDefaultPullRequestFileList(),
		MockBuiltIns(),
		nil,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.EqualError(t, err, "error getting issues mock response not found for /repos/foobar/default-mock-repo/issues/6/comments")
}

func TestReport_OnSilentMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var isDeletedCommentRequested bool
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String(fmt.Sprintf("%s\n**Reviewpad Report**\n\n:scroll: **Explanation**\nNo workflows activated", ReviewpadReportCommentAnnotation)),
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
		nil,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.True(t, isDeletedCommentRequested)
}

func TestReport_OnSilentMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
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
		nil,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.SILENT_MODE, false)

	assert.Nil(t, err)
	assert.False(t, isDeletedCommentRequested)
}

func TestReport_OnVerboseMode_WhenNoReviewpadCommentIsFound(t *testing.T) {
	var addedComment string
	commentToBeAdded := fmt.Sprintf("%s\n**Reviewpad Report**\n\n:scroll: **Executed actions**\n```yaml\nNo actions executed\n```\n", ReviewpadReportCommentAnnotation)
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
					rawBody, _ := io.ReadAll(r.Body)
					body := github.IssueComment{}

					utils.MustUnmarshal(rawBody, &body)

					addedComment = *body.Body
				}),
			),
		},
		nil,
		MockBuiltIns(),
		nil,
	)

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	err := mockedInterpreter.Report(engine.VERBOSE_MODE, false)

	assert.Nil(t, err)
	assert.Equal(t, commentToBeAdded, addedComment)
}

func TestReport_OnVerboseMode_WhenThereIsAlreadyAReviewpadComment(t *testing.T) {
	var updatedComment string
	commentUpdated := fmt.Sprintf("%s\n**Reviewpad Report**\n\n:scroll: **Executed actions**\n```yaml\nNo actions executed\n```\n", ReviewpadReportCommentAnnotation)
	mockedEnv := MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
				[]*github.IssueComment{
					{
						ID:   github.Int64(1234),
						Body: github.String(fmt.Sprintf("%s\n**Reviewpad Report**\n**:information_source: Messages**\n* Changes the README.md", ReviewpadReportCommentAnnotation)),
					},
				},
			),
			mock.WithRequestMatchHandler(
				mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					rawBody, _ := io.ReadAll(r.Body)
					body := github.IssueComment{}

					utils.MustUnmarshal(rawBody, &body)

					updatedComment = *body.Body
				}),
			),
		},
		nil,
		MockBuiltIns(),
		nil,
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
	mockErr := errors.New("mock error")

	codehostClient := GetDefaultCodeHostClient(t, nil, nil, mockErr, nil)

	// TODO: Ideally, we should not have nil arguments in the call to NewInterpreter
	gotInterpreter, err := NewInterpreter(
		ctx,
		DefaultMockLogger,
		false,
		nil,
		codehostClient,
		nil,
		DefaultMockTargetEntity,
		nil,
		nil,
		nil,
	)

	assert.Nil(t, gotInterpreter)
	assert.NotNil(t, err)
}

func TestNewInterpreter(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	wantInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	gotInterpreter, err := NewInterpreter(
		mockedEnv.GetCtx(),
		DefaultMockLogger,
		mockedEnv.GetDryRun(),
		mockedEnv.GetGithubClient(),
		mockedEnv.GetCodeHostClient(),
		mockedEnv.GetCollector(),
		DefaultMockTargetEntity,
		mockedEnv.GetEventPayload(),
		mockedEnv.GetBuiltIns(),
		nil,
	)

	assert.Nil(t, err)
	assert.Equal(t, wantInterpreter, gotInterpreter)
}

func TestReportMetric(t *testing.T) {
	commentCreated := false
	commentUpdated := false
	tests := map[string]struct {
		clientOptions          []mock.MockBackendOption
		graphQLHandler         func(res http.ResponseWriter, req *http.Request)
		commentShouldBeCreated bool
		commentShouldBeUpdated bool
		err                    error
	}{
		"when getting first commit and review date failed": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(http.StatusBadRequest)
			},
			err: errors.New("non-200 OK status code: 400 Bad Request body: \"\""),
		},
		"when find report comment failed": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				utils.MustWrite(
					res,
					`{
						"data": {
							"repository": {
								"pullRequest": {
									"commits": {
										"nodes": [
											{
												"commit": {
													"authoredDate": "2022-10-26T07:53:27Z"
												}
											}
										]
									},
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`,
				)
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusExpectationFailed)
					}),
				),
			},
			err: errors.New("error getting issues "),
		},
		"when create comment failed": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				utils.MustWrite(
					res,
					`{
						"data": {
							"repository": {
								"pullRequest": {
									"commits": {
										"nodes": [
											{
												"commit": {
													"authoredDate": "2022-10-26T07:53:27Z"
												}
											}
										]
									},
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`,
				)
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}),
				),
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					[]*github.IssueComment{},
				),
			},
			err: errors.New("error on creating report comment "),
		},
		"when update comment failed": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				utils.MustWrite(
					res,
					`{
						"data": {
							"repository": {
								"pullRequest": {
									"commits": {
										"nodes": [
											{
												"commit": {
													"authoredDate": "2022-10-26T07:53:27Z"
												}
											}
										]
									},
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`,
				)
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
					}),
				),
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					[]*github.IssueComment{
						{
							ID:   github.Int64(1),
							Body: github.String(ReviewpadMetricReportCommentAnnotation),
						},
					},
				),
			},
			err: errors.New("error on updating report comment "),
		},
		"when successfully created report comment": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				utils.MustWrite(
					res,
					`{
						"data": {
							"repository": {
								"pullRequest": {
									"commits": {
										"nodes": [
											{
												"commit": {
													"authoredDate": "2022-10-26T07:53:27Z"
												}
											}
										]
									},
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`,
				)
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						commentCreated = true
					}),
				),
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					[]*github.IssueComment{},
				),
			},
			commentShouldBeCreated: true,
			err:                    nil,
		},
		"when successfully updated report comment": {
			graphQLHandler: func(res http.ResponseWriter, req *http.Request) {
				utils.MustWrite(
					res,
					`{
						"data": {
							"repository": {
								"pullRequest": {
									"commits": {
										"nodes": [
											{
												"commit": {
													"authoredDate": "2022-10-26T07:53:27Z"
												}
											}
										]
									},
									"reviews": {
										"nodes": []
									}
								}
							}
						}
					}`,
				)
			},
			clientOptions: []mock.MockBackendOption{
				mock.WithRequestMatchHandler(
					mock.PatchReposIssuesCommentsByOwnerByRepoByCommentId,
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						commentUpdated = true
					}),
				),
				mock.WithRequestMatch(
					mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber,
					[]*github.IssueComment{
						{
							ID:   github.Int64(1),
							Body: github.String(ReviewpadMetricReportCommentAnnotation),
						},
					},
				),
			},
			commentShouldBeUpdated: true,
			err:                    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := MockDefaultEnv(t, test.clientOptions, test.graphQLHandler, nil, nil)
			interpreter, err := NewInterpreter(
				env.GetCtx(),
				env.GetLogger(),
				env.GetDryRun(),
				env.GetGithubClient(),
				env.GetCodeHostClient(),
				env.GetCollector(),
				env.GetTarget().GetTargetEntity(),
				nil,
				nil,
				nil,
			)

			assert.Nil(t, err)

			err = interpreter.ReportMetrics()

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.commentShouldBeCreated, commentCreated)
			assert.Equal(t, test.commentShouldBeUpdated, commentUpdated)

			commentCreated = false
			commentUpdated = false
		})
	}
}
