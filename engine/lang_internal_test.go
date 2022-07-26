// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals_WhenPadImportsAreEqual(t *testing.T) {
	padImport := PadImport{"http://foo.bar"}
	otherPadImport := PadImport{"http://foo.bar"}
	assert.True(t, padImport.equals(otherPadImport))
}

func TestEquals_WhenPadImportsAreDiff(t *testing.T) {
	padImport := PadImport{"http://foo.bar1"}
	otherPadImport := PadImport{"http://foo.bar2"}
	assert.False(t, padImport.equals(otherPadImport))
}

func TestEquals_WhenPadRulesAreEqual(t *testing.T) {
	padRule := PadRule{
		Name:        "test-rule",
		Kind:        "patch",
		Description: "testing rule",
		Spec:        "1 == 1",
	}

	otherPadRule := PadRule{
		Name:        "test-rule",
		Kind:        "patch",
		Description: "testing rule",
		Spec:        "1 == 1",
	}

	assert.True(t, padRule.equals(otherPadRule))
}

func TestEquals_WhenPadRulesAreDiff(t *testing.T) {
	padRule := PadRule{
		Name:        "test-rule-1",
		Kind:        "patch",
		Description: "testing rule #1",
		Spec:        "1 == 1",
	}

	otherPadRule := PadRule{
		Name:        "test-rule-2",
		Kind:        "patch",
		Description: "testing rule #2",
		Spec:        "1 < 2",
	}

	assert.False(t, padRule.equals(otherPadRule))
}

func TestEquals_WhenPadWorkflowRulesAreEqual(t *testing.T) {
	padWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction()",
		},
	}

	otherPadWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction()",
		},
	}

	assert.True(t, padWorkflowRule.equals(otherPadWorkflowRule))
}

func TestEquals_WhenPadWorkflowRulesHaveDiffRules(t *testing.T) {
	padWorkflowRule := PadWorkflowRule{
		Rule: "test-rule-1",
		ExtraActions: []string{
			"$extraAction()",
		},
	}

	otherPadWorkflowRule := PadWorkflowRule{
		Rule: "test-rule-2",
		ExtraActions: []string{
			"$extraAction()",
		},
	}

	assert.False(t, padWorkflowRule.equals(otherPadWorkflowRule))
}

func TestEquals_WhenPadWorkflowRulesHaveDiffExtraActionsLength(t *testing.T) {
	padWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction1()",
		},
	}

	otherPadWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction1()",
			"$extraAction2()",
		},
	}

	assert.False(t, padWorkflowRule.equals(otherPadWorkflowRule))
}

func TestEquals_WhenPadWorkflowRulesHaveDiffExtraActions(t *testing.T) {
	padWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction1()",
		},
	}

	otherPadWorkflowRule := PadWorkflowRule{
		Rule: "test-rule",
		ExtraActions: []string{
			"$extraAction2()",
		},
	}

	assert.False(t, padWorkflowRule.equals(otherPadWorkflowRule))
}

func TestEquals_WhenPadLabelsAreEqual(t *testing.T) {
	padLabel := PadLabel{
		Name:        "bug",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	otherPadLabel := PadLabel{
		Name:        "bug",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	assert.True(t, padLabel.equals(otherPadLabel))
}

func TestEquals_WhenPadLabelsHaveDiffName(t *testing.T) {
	padLabel := PadLabel{
		Name:        "bug#1",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	otherPadLabel := PadLabel{
		Name:        "bug#2",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	assert.False(t, padLabel.equals(otherPadLabel))
}

func TestEquals_WhenPadLabelsHaveDiffColor(t *testing.T) {
	padLabel := PadLabel{
		Name:        "bug",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	otherPadLabel := PadLabel{
		Name:        "bug",
		Color:       "a2eeef",
		Description: "Something isn't working",
	}

	assert.False(t, padLabel.equals(otherPadLabel))
}

func TestEquals_WhenPadLabelsHaveDiffDescription(t *testing.T) {
	padLabel := PadLabel{
		Name:        "bug",
		Color:       "f29513",
		Description: "Something isn't working",
	}

	otherPadLabel := PadLabel{
		Name:        "bug",
		Color:       "a2eeef",
		Description: "Something isn't working x2",
	}

	assert.False(t, padLabel.equals(otherPadLabel))
}

func TestEquals_WhenPadWorkflowAreEqual(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.True(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowHaveDiffName(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test#1",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test#2",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffDescription(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process #1",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process #2",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffRulesLength(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
			{
				Rule:         "complexChange",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffRules(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology#1",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology#2",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffActionsLength(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action1()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action1()",
			"$action2()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffActions(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action1()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action2()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadWorkflowsHaveDiffAlwaysRun(t *testing.T) {
	padWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   true,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	otherPadWorkflow := PadWorkflow{
		Name:        "test",
		Description: "Test process",
		AlwaysRun:   false,
		Rules: []PadWorkflowRule{
			{
				Rule:         "tautology",
				ExtraActions: []string{},
			},
		},
		Actions: []string{
			"$action()",
		},
	}

	assert.False(t, padWorkflow.equals(otherPadWorkflow))
}

func TestEquals_WhenPadGroupsAreEqual(t *testing.T) {
	padGroup := PadGroup{
		Name:        "juniors",
		Description: "Group of junior developers",
		Kind:        "developers",
		Type:        "filter",
		Param:       "dev",
		Where:       "$totalCreatedPullRequests($dev) < 10",
	}

	otherPadGroup := PadGroup{
		Name:        "juniors",
		Description: "Group of junior developers",
		Kind:        "developers",
		Type:        "filter",
		Param:       "dev",
		Where:       "$totalCreatedPullRequests($dev) < 10",
	}

	assert.True(t, padGroup.equals(otherPadGroup))
}

func TestEquals_WhenPadGroupsDiff(t *testing.T) {
	padGroup := PadGroup{
		Name:        "juniors",
		Description: "Group of junior developers",
		Kind:        "developers",
		Type:        "filter",
		Param:       "dev",
		Where:       "$totalCreatedPullRequests($dev) < 10",
	}

	otherPadGroup := PadGroup{
		Name:        "seniors",
		Description: "Senior developers",
		Kind:        "developers",
		Spec:        "[\"john\"]",
	}

	assert.False(t, padGroup.equals(otherPadGroup))
}

func TestEquals_WhenReviewpadFilesAreEqual(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.True(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffVersion(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1beta",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffEdition(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffIgnoreErrors(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: true,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfImports(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
			{Url: "https://foo.bar/tautology-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffImports(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/tautology-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfRules(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
			{
				Name:        "test-rule-2",
				Kind:        "patch",
				Description: "testing rule #2",
				Spec:        "1 < 2",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffRules(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule-2",
				Kind:        "patch",
				Description: "testing rule #2",
				Spec:        "1 < 2",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfLabels(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
			"bug#2": {
				Name:        "bug#2",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffLabels(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug#2": {
				Name:        "bug#2",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfWorkflows(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
					{
						Rule:         "tautology",
						ExtraActions: []string{},
					},
				},
				Actions: []string{
					"$action()",
				},
			},
			{
				Name:        "test#2",
				Description: "Test process x2",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
					{
						Rule:         "test-rule",
						ExtraActions: []string{},
					},
				},
				Actions: []string{
					"$action2()",
				},
			},
		},
	}

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffWorkflows(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test#2",
				Description: "Test process x2",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
					{
						Rule:         "test-rule",
						ExtraActions: []string{},
					},
				},
				Actions: []string{
					"$action2()",
				},
			},
		},
	}

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfGroups(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
			{
				Name:        "juniors",
				Description: "Group of junior developers",
				Kind:        "developers",
				Type:        "filter",
				Param:       "dev",
				Where:       "$totalCreatedPullRequests($dev) < 10",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffGroups(t *testing.T) {
	reviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "seniors",
				Description: "Senior developers",
				Kind:        "developers",
				Spec:        "[\"john\"]",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	otherReviewpadFile := &ReviewpadFile{
		Version:      "reviewpad.com/v1alpha",
		Edition:      "professional",
		Mode:         "silent",
		IgnoreErrors: false,
		Imports: []PadImport{
			{Url: "https://foo.bar/draft-rule.yml"},
		},
		Groups: []PadGroup{
			{
				Name:        "juniors",
				Description: "Group of junior developers",
				Kind:        "developers",
				Type:        "filter",
				Param:       "dev",
				Where:       "$totalCreatedPullRequests($dev) < 10",
			},
		},
		Rules: []PadRule{
			{
				Name:        "test-rule",
				Kind:        "patch",
				Description: "testing rule",
				Spec:        "1 == 1",
			},
		},
		Labels: map[string]PadLabel{
			"bug": {
				Name:        "bug",
				Color:       "f29513",
				Description: "Something isn't working",
			},
		},
		Workflows: []PadWorkflow{
			{
				Name:        "test",
				Description: "Test process",
				AlwaysRun:   true,
				Rules: []PadWorkflowRule{
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

	assert.False(t, reviewpadFile.equals(otherReviewpadFile))
}
