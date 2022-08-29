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

	initialStat := BuildStatement("$actionB()")

	programUnderTest := BuildProgram([]*Statement{initialStat})

	addedStat := BuildStatement(action)
	wantProgram := BuildProgram([]*Statement{
		initialStat,
		addedStat,
	})

	programUnderTest.append(workflow.Actions)

	assert.Equal(t, wantProgram, programUnderTest)
}
