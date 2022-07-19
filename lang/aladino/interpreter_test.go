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

func TestExecProgram_WhenExecStatementFails(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
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

	assert.EqualError(t, err, "no type for built-in action. Please check if the mode in the reviewpad.yml file supports it.")
}

func TestExecProgram(t *testing.T) {
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

	program := &engine.Program{
		Statements: []*engine.Statement{
			statement,
		},
	}

	err = mockedInterpreter.ExecProgram(program)

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
