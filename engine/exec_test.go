// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/jinzhu/copier"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

// TODO: Move this to resources/test/engine once #174 is accepted
var mockedReviewpadFile = &engine.ReviewpadFile{
	Version:      "reviewpad.com/v1alpha",
	Edition:      "professional",
	Mode:         "silent",
	IgnoreErrors: false,
	Imports: []engine.PadImport{
		{Url: "https://foo.bar/draft-rule.yml"},
	},
	Groups: []engine.PadGroup{
		{
			Name:        "seniors",
			Description: "Senior developers",
			Kind:        "developers",
			Spec:        "[\"john\"]",
		},
	},
	Rules: []engine.PadRule{
		{
			Name:        "tautology",
			Kind:        "patch",
			Description: "testing rule",
			Spec:        "true",
		},
	},
	Labels: map[string]engine.PadLabel{
		"bug": {
			Name:        "bug",
			Color:       "f29513",
			Description: "Something isn't working",
		},
	},
	Workflows: []engine.PadWorkflow{
		{
			Name:        "test-workflow",
			Description: "Test process",
			AlwaysRun:   true,
			Rules: []engine.PadWorkflowRule{
				{
					Rule:         "tautology",
					ExtraActions: []string{},
				},
			},
			Actions: []string{
				"$action()",
			},
		},
	},
}

func TestEval_WhenDryModeIsNotSetAndGetLabelRequestFails(t *testing.T) {
	failMessage := "GetLabelRequestFailed"
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposLabelsByOwnerByRepoByName,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(mock.MustMarshal(github.ErrorResponse{
						// An error response may also consist of a 404 status code.
						// However, in this context, such response is not an error since it only means a label does not exist.
						// So we force an error by setting it to a status code different than 404.
						Response: &http.Response{
							StatusCode: 500,
						},
						Message: failMessage,
					}))
				}),
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(mockedReviewpadFile, mockedEnv)

	assert.Nil(t, gotProgram)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestEval_WhenDryModeIsNotSetAndCreateLabelRequestFails(t *testing.T) {
	failMessage := "CreateLabelRequestFailed"
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposLabelsByOwnerByRepoByName,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(mock.MustMarshal(github.ErrorResponse{
						Response: &http.Response{
							StatusCode: 404,
						},
						Message: "Resource not found",
					}))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.PostReposLabelsByOwnerByRepo,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(mockedReviewpadFile, mockedEnv)

	assert.Nil(t, gotProgram)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestEval_WhenWorkflowRuleEvalExprFails(t *testing.T) {
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{
					Name: github.String("bug"),
				},
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	otherReviewpadFile := &engine.ReviewpadFile{}
	copier.Copy(otherReviewpadFile, mockedReviewpadFile)

	executedWorkflow := engine.PadWorkflow{
		Name:        "test-workflow",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []engine.PadWorkflowRule{
			{
				Rule:         "wrong-typed-rule",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	wrongTypedRule := engine.PadRule{
		Name:        "wrong-typed-rule",
		Kind:        "patch",
		Description: "Rule with non boolean spec",
		Spec:        "\"notBoolType\"",
	}

	otherReviewpadFile.Rules = []engine.PadRule{wrongTypedRule}
	otherReviewpadFile.Workflows = []engine.PadWorkflow{executedWorkflow}

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(otherReviewpadFile, mockedEnv)

	assert.Nil(t, gotProgram)
	assert.EqualError(t, err, "expression \"notBoolType\" is not a condition")
}

func TestEval_WhenWorkflowIsTriggered(t *testing.T) {
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{
					Name: github.String("bug"),
				},
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	otherReviewpadFile := &engine.ReviewpadFile{}
	copier.Copy(otherReviewpadFile, mockedReviewpadFile)

	executedWorkflow := engine.PadWorkflow{
		Name:        "test-workflow",
		Description: "Test process",
		AlwaysRun:   false,
		Rules: []engine.PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherReviewpadFile.Workflows = []engine.PadWorkflow{executedWorkflow}

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(otherReviewpadFile, mockedEnv)

	wantProgram := &engine.Program{
		Statements: []*engine.Statement{
			{
				Code: executedWorkflow.Actions[0],
				Metadata: &engine.Metadata{
					Workflow:    executedWorkflow,
					TriggeredBy: executedWorkflow.Rules,
				},
			},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantProgram, gotProgram)
}

func TestEval_WhenWorkflowIsSkipped(t *testing.T) {
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{
					Name: github.String("bug"),
				},
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	otherReviewpadFile := &engine.ReviewpadFile{}
	copier.Copy(otherReviewpadFile, mockedReviewpadFile)

	executedWorkflow := engine.PadWorkflow{
		Name:        "test-workflow-A",
		Description: "Test process A",
		AlwaysRun:   false,
		Rules: []engine.PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$actionA()",
		},
	}

	rule := engine.PadRule{
		Name:        "dummy-rule",
		Kind:        "patch",
		Description: "testing rule",
		Spec:        "true",
	}

	otherReviewpadFile.Rules = append(otherReviewpadFile.Rules, rule)

	otherReviewpadFile.Workflows = []engine.PadWorkflow{
		executedWorkflow,
		{
			Name:        "test-workflow-B",
			Description: "Test process B",
			AlwaysRun:   false,
			Rules: []engine.PadWorkflowRule{
				{
					Rule:         "dummy-rule",
					ExtraActions: []string{},
				},
			},
			Actions: []string{
				"$actionB()",
			},
		},
	}

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(otherReviewpadFile, mockedEnv)

	wantProgram := &engine.Program{
		Statements: []*engine.Statement{
			{
				Code: executedWorkflow.Actions[0],
				Metadata: &engine.Metadata{
					Workflow:    executedWorkflow,
					TriggeredBy: executedWorkflow.Rules,
				},
			},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantProgram, gotProgram)
}

func TestEval_WhenNoWorkflowRulesAreActivated(t *testing.T) {
	mockedClient := engine.MockGithubClient(
		[]mock.MockBackendOption{
			mock.WithRequestMatch(
				mock.GetReposLabelsByOwnerByRepoByName,
				&github.Label{
					Name: github.String("bug"),
				},
			),
		},
	)

	dryRun := false
	mockedCtx := engine.DefaultMockCtx
	mockedCollector := engine.DefaultMockCollector
	mockedEvent := engine.DefaultMockEventPayload
	mockedPullRequest := engine.GetDefaultMockPullRequestDetails()

	otherReviewpadFile := &engine.ReviewpadFile{}
	copier.Copy(otherReviewpadFile, mockedReviewpadFile)

	notTriggeredRule := engine.PadRule{
		Name:        "dummy-rule",
		Kind:        "patch",
		Description: "testing rule",
		Spec:        "false",
	}

	analyzedWorkflow := engine.PadWorkflow{
		Name:        "test-workflow",
		Description: "Test process",
		AlwaysRun:   false,
		Rules: []engine.PadWorkflowRule{
			{
				Rule:         "dummy-rule",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherReviewpadFile.Rules = []engine.PadRule{notTriggeredRule}
	otherReviewpadFile.Workflows = []engine.PadWorkflow{analyzedWorkflow}

	mockedAladinoInterpreter, err := aladino.NewInterpreter(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		aladino.MockBuiltIns(),
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("aladino NewInterpreter returned unexpected error: %v", err))
	}

	mockedEnv, err := engine.NewEvalEnv(
		mockedCtx,
		dryRun,
		mockedClient,
		// TODO: add mocked github GQL client
		nil,
		mockedCollector,
		mockedPullRequest,
		mockedEvent,
		mockedAladinoInterpreter,
	)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("NewEvalEnv returned unexpected error: %v", err))
	}

	gotProgram, err := engine.Eval(otherReviewpadFile, mockedEnv)

	wantProgram := &engine.Program{Statements: []*engine.Statement{}}

	assert.Nil(t, err)
	assert.Equal(t, wantProgram, gotProgram)
}
