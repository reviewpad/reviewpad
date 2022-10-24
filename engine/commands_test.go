// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/stretchr/testify/assert"
)

func TestAssignReviewer(t *testing.T) {
	testCases := map[string]struct {
		matches  []string
		rule     *engine.PadRule
		workflow *engine.PadWorkflow
		err      error
	}{
		"when arguments are empty": {
			matches: []string{},
			err:     errors.New("invalid assign reviewer command"),
		},
		"when invalid number of arguments": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa, shay2025 2",
			},
			err: errors.New("invalid assign reviewer command"),
		},
		"when number of reviewers is not a number": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa, shay2025 z random",
				"marcelosousa, shay2025",
				"z",
				"random",
			},
			err: &strconv.NumError{Func: "Atoi", Num: "z", Err: errors.New("invalid syntax")},
		},
		"when missing number of reviewers and strategy": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa",
				"marcelosousa",
			},
			rule: &engine.PadRule{
				Spec: "true",
			},
			workflow: &engine.PadWorkflow{
				On:          []handler.TargetEntityKind{handler.PullRequest},
				Description: "assign reviewer command",
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "assign-reviewer",
					},
				},
				AlwaysRun: true,
				Actions:   []string{`$assignReviewer(["marcelosousa"], 1, "reviewpad")`},
			},
		},
		"when missing number of strategy": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa, shay2025 1",
				"marcelosousa, shay2025",
				"1",
			},
			rule: &engine.PadRule{
				Spec: "true",
			},
			workflow: &engine.PadWorkflow{
				On:          []handler.TargetEntityKind{handler.PullRequest},
				Description: "assign reviewer command",
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "assign-reviewer",
					},
				},
				AlwaysRun: true,
				Actions:   []string{`$assignReviewer(["marcelosousa","shay2025"], 1, "reviewpad")`},
			},
		},
		"when one reviewer": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa 1 reviewpad",
				"marcelosousa",
				"1",
				"reviewpad",
			},
			rule: &engine.PadRule{
				Spec: "true",
			},
			workflow: &engine.PadWorkflow{
				On:          []handler.TargetEntityKind{handler.PullRequest},
				Description: "assign reviewer command",
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "assign-reviewer",
					},
				},
				AlwaysRun: true,
				Actions:   []string{`$assignReviewer(["marcelosousa"], 1, "reviewpad")`},
			},
		},
		"when two reviewer": {
			matches: []string{
				"/reviewpad assign-reviewer marcelosousa, shay2025 2 random",
				"marcelosousa, shay2025",
				"2",
				"random",
			},
			rule: &engine.PadRule{
				Spec: "true",
			},
			workflow: &engine.PadWorkflow{
				On:          []handler.TargetEntityKind{handler.PullRequest},
				Description: "assign reviewer command",
				Rules: []engine.PadWorkflowRule{
					{
						Rule: "assign-reviewer",
					},
				},
				AlwaysRun: true,
				Actions:   []string{`$assignReviewer(["marcelosousa","shay2025"], 2, "random")`},
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			rule, workflow, err := engine.AssignReviewerCommand(test.matches)

			// since we are attaching a random string
			// to the dynamically generated rule name
			// we are setting them to the same value
			if rule != nil {
				test.rule.Name = rule.Name
				test.workflow.Name = rule.Name
				test.workflow.Rules[0].Rule = rule.Name
			}

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.rule, rule)
			assert.Equal(t, test.workflow, workflow)
		})
	}
}
