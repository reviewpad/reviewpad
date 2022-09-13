// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func inlineRulesNormalizer() *NormalizeRule {
	normalizedRule := NewNormalizeRule()
	normalizedRule.WithModificators(inlineRulesModificator)
	return normalizedRule
}

func inlineRulesModificator(file *ReviewpadFile) (*ReviewpadFile, error) {
	reviewpadFile := &ReviewpadFile{
		Version:      file.Version,
		Edition:      file.Edition,
		Mode:         file.Mode,
		IgnoreErrors: file.IgnoreErrors,
		Imports:      file.Imports,
		Groups:       file.Groups,
		Rules:        file.Rules,
		Labels:       file.Labels,
		Workflows:    file.Workflows,
	}

	for i, workflow := range reviewpadFile.Workflows {
		processedWorkflow, rules, err := processInlineRulesOnWorkflow(workflow, reviewpadFile.Rules)
		if err != nil {
			return nil, err
		}

		reviewpadFile.Rules = append(reviewpadFile.Rules, rules...)
		reviewpadFile.Workflows[i] = *processedWorkflow
	}

	return reviewpadFile, nil
}

func processInlineRulesOnWorkflow(workflow PadWorkflow, currentRules []PadRule) (*PadWorkflow, []PadRule, error) {
	wf := &PadWorkflow{
		Name:        workflow.Name,
		Description: workflow.Description,
		AlwaysRun:   workflow.AlwaysRun,
		Rules:       workflow.Rules,
		Actions:     workflow.Actions,
		On:          workflow.On,
	}
	rules := make([]PadRule, 0)

	for _, rawRule := range workflow.NonNormalizedRules {
		var rule *PadRule
		var workflowRule *PadWorkflowRule

		switch r := rawRule.(type) {
		case string:
			rule = decodeRule(r)
			workflowRule = &PadWorkflowRule{
				Rule:         rule.Name,
				ExtraActions: []string{},
			}
		case map[string]interface{}:
			decodedWorkflowRule, err := decodeWorkflowRule(r)
			if err != nil {
				return nil, nil, err
			}
			workflowRule = decodedWorkflowRule
			rule = decodeRule(decodedWorkflowRule.Rule)
		default:
			return nil, nil, fmt.Errorf("unknown rule type %T", r)
		}

		if _, exists := findRule(currentRules, rule.Name); !exists {
			rules = append(rules, *rule)
		}

		if workflowRule != nil {
			wf.Rules = append(wf.Rules, *workflowRule)
		}
	}

	workflow.NonNormalizedRules = nil

	return wf, rules, nil
}

func decodeRule(rule string) *PadRule {
	return &PadRule{
		Name: rule,
		Spec: rule,
		Kind: "patch",
	}
}

func decodeWorkflowRule(rule map[string]interface{}) (*PadWorkflowRule, error) {
	workflowRule := &PadWorkflowRule{}
	err := mapstructure.Decode(rule, workflowRule)
	return workflowRule, err
}
