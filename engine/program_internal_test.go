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

	initialStatMetadata := BuildMetadata(
		PadWorkflow{
			Name: "test-workflow-B",
		},
		[]PadWorkflowRule{
			{Rule: "test-rule-B"},
		},
	)
	initialStat := BuildStatement("$actionB()", initialStatMetadata)

	programUnderTest := BuildProgram([]*Statement{initialStat})

	addedStatMetadata := BuildMetadata(workflow, workflow.Rules)
	addedStat := BuildStatement(action, addedStatMetadata)
	wantProgram := BuildProgram([]*Statement{
		initialStat,
		addedStat,
	})

	programUnderTest.append(workflow.Actions, workflow, workflow.Rules)

	assert.Equal(t, wantProgram, programUnderTest)
}
