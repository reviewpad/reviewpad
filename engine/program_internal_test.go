// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppend(t *testing.T) {
	workflow := PadWorkflow{
		Name:        "test-workflow-A",
		Description: "Testing workflow",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{Rule: "test-rule-A"},
		},
		Actions: []string{"$actionA()"},
	}
	actions := []string{"$actionB()"}
	rules := []PadWorkflowRule{
		{Rule: "test-rule-B"},
	}

	statement := &Statement{
		Code: "$actionC()",
		Metadata: &Metadata{
			Workflow: PadWorkflow{
				Name: "test-workflow-C",
			},
			TriggeredBy: []PadWorkflowRule{
				{Rule: "test-rule-C"},
			},
		},
	}
	program := &Program{
		Statements: []*Statement{statement},
	}

	wantProgram := &Program{
		Statements: []*Statement{
			statement,
			{
				Code: "$actionB()",
				Metadata: &Metadata{
					Workflow:    workflow,
					TriggeredBy: rules,
				},
			},
		},
	}

	program.append(actions, workflow, rules)

	assert.Equal(t, wantProgram, program)
}
