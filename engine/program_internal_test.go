// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	action := "$actionA()"
	workflow := PadWorkflow{
		Name:        "test-workflow-A",
		Description: "Testing workflow",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{Rule: "test-rule-A"},
		},
		Actions: []string{action},
	}

	initialStatement := &Statement{
		Code: "$actionB()",
		Metadata: &Metadata{
			Workflow: PadWorkflow{
				Name: "test-workflow-B",
			},
			TriggeredBy: []PadWorkflowRule{
				{Rule: "test-rule-B"},
			},
		},
	}

	programUnderTest := &Program{
		Statements: []*Statement{initialStatement},
	}

	wantProgram := &Program{
		Statements: []*Statement{
			initialStatement,
			{
				Code: action,
				Metadata: &Metadata{
					Workflow:    workflow,
					TriggeredBy: workflow.Rules,
				},
			},
		},
	}

	programUnderTest.append(workflow.Actions, workflow, workflow.Rules)

	assert.Equal(t, wantProgram, programUnderTest)
}
