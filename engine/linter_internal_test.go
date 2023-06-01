// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallsToRuleBuiltIn(t *testing.T) {
	wantRuleNames := []string{
		"rule1",
		"rule2",
		"rule3",
		"rule4",
		"rule-five",
	}

	groups := []PadGroup{}
	rules := []PadRule{
		{
			Spec: `$rule("rule1") && $rule("rule2")`,
		},
		{
			Spec: `true`,
		},
		{
			Spec: `$rule("rule3") && $test()`,
		},
		{
			Spec: `$test() && $rule("rule4")`,
		},
		{
			Spec: `$test1()`,
		},
		{
			Spec: `$rule("rule-five") && $isElementOf($author(), $group("group-name"))`,
		},
	}
	workflows := []PadWorkflow{}

	gotRuleNames := getCallsToRuleBuiltIn(groups, rules, workflows)

	assert.Equal(t, wantRuleNames, gotRuleNames)
}

func TestShadowedVariable(t *testing.T) {
	workflow := PadWorkflow{
		Runs: []PadWorkflowRunBlock{
			{
				ForEach: &PadWorkflowRunForEachBlock{
					Value: "var1",
					Do: []PadWorkflowRunBlock{
						{
							ForEach: &PadWorkflowRunForEachBlock{
								Value: "var2",
							},
						},
					},
				},
			},
			{
				ForEach: &PadWorkflowRunForEachBlock{
					Value: "var1",
					Do: []PadWorkflowRunBlock{
						{
							ForEach: &PadWorkflowRunForEachBlock{
								Value: "var1",
							},
						},
					},
				},
			},
		},
	}

	wantErr := errors.New("variable shadowing is not allowed: the variable `var1` is already defined")

	gotErr := lintShadowedVariables([]PadWorkflow{workflow})

	assert.Equal(t, wantErr, gotErr)
}
