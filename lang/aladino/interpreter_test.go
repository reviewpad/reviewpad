// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"log"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v2/engine"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	"github.com/stretchr/testify/assert"
)

func TestExecStatement_WhenParseFails(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
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
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
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
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
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
	mockedEnv, err := mocks_aladino.MockDefaultEnv(
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
	)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
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

	wantVal := aladino.ReportWorkflowDetails{
		Name: statementWorkflowName,
		Rules: map[string]bool{
			statementRule: true,
		},
		Actions: []string{statementCode},
	}

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
