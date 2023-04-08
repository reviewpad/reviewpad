// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func inlineNormalizer() *NormalizeRule {
	normalizedRule := NewNormalizeRule()
	normalizedRule.WithModificators(inlineModificator)
	return normalizedRule
}

func inlineModificator(file *ReviewpadFile) (*ReviewpadFile, error) {
	reviewpadFile := &ReviewpadFile{
		Mode:           file.Mode,
		IgnoreErrors:   file.IgnoreErrors,
		MetricsOnMerge: file.MetricsOnMerge,
		Extends:        file.Extends,
		Imports:        file.Imports,
		Groups:         file.Groups,
		Rules:          file.Rules,
		Labels:         file.Labels,
		Workflows:      file.Workflows,
		Pipelines:      file.Pipelines,
		Recipes:        file.Recipes,
	}

	for i, workflow := range reviewpadFile.Workflows {
		processedWorkflow, rules, err := processWorkflow(workflow, reviewpadFile.Rules)
		if err != nil {
			return nil, err
		}

		reviewpadFile.Rules = append(reviewpadFile.Rules, rules...)
		reviewpadFile.Workflows[i] = *processedWorkflow
	}

	for _, pipeline := range reviewpadFile.Pipelines {
		for _, stage := range pipeline.Stages {
			actions, err := processCompactActions(stage.NonNormalizedActions)
			if err != nil {
				return nil, err
			}

			stage.Actions = actions
			stage.NonNormalizedActions = nil
		}
	}

	return reviewpadFile, nil
}

func processWorkflow(workflow PadWorkflow, currentRules []PadRule) (*PadWorkflow, []PadRule, error) {
	wf := &PadWorkflow{
		Name:        workflow.Name,
		Description: workflow.Description,
		AlwaysRun:   workflow.AlwaysRun,
		Rules:       workflow.Rules,
		Actions:     workflow.Actions,
		On:          workflow.On,
	}

	actions, err := processCompactActions(workflow.NonNormalizedActions)
	if err != nil {
		return nil, nil, err
	}

	wf.Actions = actions

	rules, workflowRules, err := processCompactRules(workflow.NonNormalizedRules, currentRules)
	if err != nil {
		return nil, nil, err
	}

	wf.Rules = workflowRules

	return wf, rules, nil
}

func decodeRule(rule string) *PadRule {
	return &PadRule{
		Name: rule,
		Spec: rule,
		Kind: "patch",
	}
}

func decodeWorkflowRule(rule map[string]any) (*PadWorkflowRule, error) {
	workflowRule := &PadWorkflowRule{}
	err := mapstructure.Decode(rule, workflowRule)
	return workflowRule, err
}

func processCompactRules(rawRule any, currentRules []PadRule) ([]PadRule, []PadWorkflowRule, error) {
	var rule *PadRule
	var workflowRule *PadWorkflowRule
	var rules []PadRule
	var workflowRules []PadWorkflowRule

	switch r := rawRule.(type) {
	case string:
		rule = decodeRule(r)
		workflowRule = &PadWorkflowRule{
			Rule: rule.Name,
		}
	case []any:
		for _, ru := range r {
			processedRules, processedWorkflowRules, err := processCompactRules(ru, currentRules)
			if err != nil {
				return nil, nil, err
			}

			rules = append(rules, processedRules...)
			workflowRules = append(workflowRules, processedWorkflowRules...)
		}
	case map[string]any:
		decodedWorkflowRule, err := decodeWorkflowRule(r)
		if err != nil {
			return nil, nil, err
		}

		workflowRule = decodedWorkflowRule
		rule = decodeRule(decodedWorkflowRule.Rule)
	default:
		return nil, nil, fmt.Errorf("unknown rule type %T", r)
	}

	// we need to check rule is not nil if we are processing a list of rules
	if rule != nil {
		if _, exists := findRule(currentRules, rule.Name); !exists {
			rules = append(rules, *rule)
		}
	}

	// we need to check workflowRule is not nil if we are processing a list of rules
	if workflowRule != nil {
		workflowRules = append(workflowRules, *workflowRule)
	}

	return rules, workflowRules, nil
}

func processCompactActions(nonNormalizedActions any) ([]string, error) {
	actions := []string{}

	switch action := nonNormalizedActions.(type) {
	case string:
		actions = append(actions, action)
	case []any:
		for _, rawAction := range action {
			processedActions, err := processCompactActions(rawAction)
			if err != nil {
				return nil, err
			}

			actions = append(actions, processedActions...)
		}
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown action type: %T", action)
	}

	return actions, nil
}
