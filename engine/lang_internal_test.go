// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockedReviewpadFile = &ReviewpadFile{
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
		Color:       "f29513",
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
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	assert.True(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffVersion(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Version = "reviewpad.com/v1beta"

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffEdition(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Edition = ""

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffMode(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Mode = "verbose"

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffIgnoreErrors(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.IgnoreErrors = true

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfImports(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Imports = []PadImport{
		{Url: "https://foo.bar/draft-rule.yml"},
		{Url: "https://foo.bar/tautology-rule.yml"},
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffImports(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Imports = []PadImport{
		{Url: "https://foo.bar/tautology-rule.yml"},
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfRules(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Rules = []PadRule{
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
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffRules(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Rules = []PadRule{
		{
			Name:        "test-rule-2",
			Kind:        "patch",
			Description: "testing rule #2",
			Spec:        "1 < 2",
		},
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfLabels(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Labels = map[string]PadLabel{
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
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffLabels(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Labels = map[string]PadLabel{
		"bug#2": {
			Name:        "bug#2",
			Color:       "f29513",
			Description: "Something isn't working",
		},
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfWorkflows(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Workflows = []PadWorkflow{
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
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffWorkflows(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Workflows = []PadWorkflow{
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
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffNumberOfGroups(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Groups = []PadGroup{
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
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestEquals_WhenReviewpadFilesHaveDiffGroups(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Groups = []PadGroup{
		{
			Name:        "juniors",
			Description: "Group of junior developers",
			Kind:        "developers",
			Type:        "filter",
			Param:       "dev",
			Where:       "$totalCreatedPullRequests($dev) < 10",
		},
	}

	assert.False(t, mockedReviewpadFile.equals(otherReviewpadFile))
}

func TestAppendLabels_WhenReviewpadFileHasNoLabels(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Labels = nil

	otherReviewpadFile.appendLabels(mockedReviewpadFile)

	wantLabels := map[string]PadLabel{
		"bug": {
			Name:        "bug",
			Color:       "f29513",
			Description: "Something isn't working",
		},
	}

	assert.Equal(t, wantLabels, otherReviewpadFile.Labels)
}

func TestAppendLabels_WhenReviewpadFileHasLabels(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Labels = map[string]PadLabel{
		"bug#2": {
			Name:        "bug#2",
			Color:       "f29513",
			Description: "Something isn't working",
		},
	}

	otherReviewpadFile.appendLabels(mockedReviewpadFile)

	wantLabels := map[string]PadLabel{
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
	}

	assert.Equal(t, wantLabels, otherReviewpadFile.Labels)
}

func TestAppendRules_WhenReviewpadFileHasNoRules(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Rules = nil

	otherReviewpadFile.appendRules(mockedReviewpadFile)

	wantRules := []PadRule{
		{
			Name:        "test-rule",
			Kind:        "patch",
			Description: "testing rule",
			Spec:        "1 == 1",
		},
	}

	assert.Equal(t, wantRules, otherReviewpadFile.Rules)
}

func TestAppendRules_WhenReviewpadFileHasRules(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Rules = []PadRule{
		{
			Name:        "test-rule-2",
			Kind:        "patch",
			Description: "testing rule #2",
			Spec:        "1 < 2",
		},
	}

	otherReviewpadFile.appendRules(mockedReviewpadFile)

	wantRules := []PadRule{
		{
			Name:        "test-rule-2",
			Kind:        "patch",
			Description: "testing rule #2",
			Spec:        "1 < 2",
		},
		{
			Name:        "test-rule",
			Kind:        "patch",
			Description: "testing rule",
			Spec:        "1 == 1",
		},
	}

	assert.Equal(t, wantRules, otherReviewpadFile.Rules)
}

func TestAppendGroups_WhenReviewpadFileHasNoGroups(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Groups = nil

	otherReviewpadFile.appendGroups(mockedReviewpadFile)

	wantGroups := []PadGroup{
		{
			Name:        "seniors",
			Description: "Senior developers",
			Kind:        "developers",
			Spec:        "[\"john\"]",
		},
	}

	assert.Equal(t, wantGroups, otherReviewpadFile.Groups)
}

func TestAppendGroups_WhenReviewpadFileHasGroups(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Groups = []PadGroup{
		{
			Name:        "juniors",
			Description: "Group of junior developers",
			Kind:        "developers",
			Type:        "filter",
			Param:       "dev",
			Where:       "$totalCreatedPullRequests($dev) < 10",
		},
	}

	otherReviewpadFile.appendGroups(mockedReviewpadFile)

	wantGroups := []PadGroup{
		{
			Name:        "juniors",
			Description: "Group of junior developers",
			Kind:        "developers",
			Type:        "filter",
			Param:       "dev",
			Where:       "$totalCreatedPullRequests($dev) < 10",
		},
		{
			Name:        "seniors",
			Description: "Senior developers",
			Kind:        "developers",
			Spec:        "[\"john\"]",
		},
	}

	assert.Equal(t, wantGroups, otherReviewpadFile.Groups)
}

func TestAppendWorkflows_WhenReviewpadFileHasNoWorkflows(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Workflows = nil

	otherReviewpadFile.appendWorkflows(mockedReviewpadFile)

	wantWorkflows := []PadWorkflow{
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
	}

	assert.Equal(t, wantWorkflows, otherReviewpadFile.Workflows)
}

func TestAppendWorkflows_WhenReviewpadFileHasWorkflows(t *testing.T) {
	otherReviewpadFile := &ReviewpadFile{}
	*otherReviewpadFile = *mockedReviewpadFile

	otherReviewpadFile.Workflows = []PadWorkflow{
		{
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
		},
	}

	otherReviewpadFile.appendWorkflows(mockedReviewpadFile)

	wantWorkflows := []PadWorkflow{
		{
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
		},
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
	}

	assert.Equal(t, wantWorkflows, otherReviewpadFile.Workflows)
}

func TestFindGroup_WhenGroupExists(t *testing.T) {
	groups := []PadGroup{
		{
			Name:        "juniors",
			Description: "Group of junior developers",
			Kind:        "developers",
			Type:        "filter",
			Param:       "dev",
			Where:       "$totalCreatedPullRequests($dev) < 10",
		},
		{
			Name:        "seniors",
			Description: "Senior developers",
			Kind:        "developers",
			Spec:        "[\"john\"]",
		},
	}

	wantGroup := &PadGroup{
		Name:        "juniors",
		Description: "Group of junior developers",
		Kind:        "developers",
		Type:        "filter",
		Param:       "dev",
		Where:       "$totalCreatedPullRequests($dev) < 10",
	}

	gotGroup, found := findGroup(groups, "juniors")

	assert.True(t, found)
	assert.Equal(t, wantGroup, gotGroup)
}

func TestFindGroup_WhenGroupDoesNotExists(t *testing.T) {
	groups := []PadGroup{
		{
			Name:        "juniors",
			Description: "Group of junior developers",
			Kind:        "developers",
			Type:        "filter",
			Param:       "dev",
			Where:       "$totalCreatedPullRequests($dev) < 10",
		},
		{
			Name:        "seniors",
			Description: "Senior developers",
			Kind:        "developers",
			Spec:        "[\"john\"]",
		},
	}

	gotGroup, found := findGroup(groups, "devs")

	assert.False(t, found)
	assert.Nil(t, gotGroup)
}

func TestFindRule_WhenRuleExists(t *testing.T) {
	rules := []PadRule{
		{
			Name:        "test-rule-1",
			Kind:        "patch",
			Description: "testing rule #1",
			Spec:        "1 == 1",
		},
		{
			Name:        "test-rule-2",
			Kind:        "patch",
			Description: "testing rule #2",
			Spec:        "1 < 2",
		},
	}

	wantRule := &PadRule{
		Name:        "test-rule-2",
		Kind:        "patch",
		Description: "testing rule #2",
		Spec:        "1 < 2",
	}

	gotRule, found := findRule(rules, "test-rule-2")

	assert.True(t, found)
	assert.Equal(t, wantRule, gotRule)
}

func TestFindRule_WhenRuleDoesNotExists(t *testing.T) {
	rules := []PadRule{
		{
			Name:        "test-rule-1",
			Kind:        "patch",
			Description: "testing rule #1",
			Spec:        "1 == 1",
		},
		{
			Name:        "test-rule-2",
			Kind:        "patch",
			Description: "testing rule #2",
			Spec:        "1 < 2",
		},
	}

	gotRule, found := findRule(rules, "test-rule-3")

	assert.False(t, found)
	assert.Nil(t, gotRule)
}
