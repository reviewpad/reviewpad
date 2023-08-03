// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
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

	wantErr := errors.New("variable shadowing is not allowed: the variable var1 is already defined")

	gotErr := lintShadowedVariables([]PadWorkflow{workflow})

	assert.Equal(t, wantErr, gotErr)
}

func TestShadowedReservedName(t *testing.T) {
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
								Value: "$comment",
							},
						},
					},
				},
			},
		},
	}

	wantErr := errors.New("built-in shadowing is not allowed: the variable $comment is a reserved name")

	gotErr := lintReservedWords([]PadWorkflow{workflow}, []string{"comment"})

	assert.Equal(t, wantErr, gotErr)
}

func TestLintRulesMention(t *testing.T) {
	tests := map[string]struct {
		rules     []PadRule
		workflows []PadWorkflow
		wantErr   error
	}{
		"when unused rule": {
			rules: []PadRule{
				{
					Name: "true == true",
					Spec: "true == true",
				},
			},
			workflows: []PadWorkflow{
				{
					Rules: []PadWorkflowRule{
						{
							Rule: "true == false",
						},
					},
				},
			},
			wantErr: nil,
		},
		"when rule used in old workflow": {
			rules: []PadRule{
				{
					Name: "true == true",
					Spec: "true == true",
				},
			},
			workflows: []PadWorkflow{
				{
					Rules: []PadWorkflowRule{
						{
							Rule: "true == true",
						},
					},
				},
			},
			wantErr: nil,
		},
		"when rule used in run": {
			rules: []PadRule{
				{
					Name: "true == true",
					Spec: "true == true",
				},
			},
			workflows: []PadWorkflow{
				{
					Runs: []PadWorkflowRunBlock{
						{
							If: []PadWorkflowRule{
								{
									Rule: "true == true",
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := lintRulesMentions(logrus.NewEntry(logrus.New()), test.rules, nil, test.workflows)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
