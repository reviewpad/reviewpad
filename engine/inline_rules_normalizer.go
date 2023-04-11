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
		for i, stage := range pipeline.Stages {
			actions, err := normalizeActions(stage.NonNormalizedActions)
			if err != nil {
				return nil, err
			}

			stage.Actions = actions
			stage.NonNormalizedActions = nil
			pipeline.Stages[i] = stage
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

	actions, err := normalizeActions(workflow.NonNormalizedActions)
	if err != nil {
		return nil, nil, err
	}

	wf.Actions = actions

	compactRules, workflowRules, err := normalizeRules(workflow.NonNormalizedRules, currentRules)
	if err != nil {
		return nil, nil, err
	}

	runs, runRules, err := processRun(workflow.NonNormalizedRun, append(currentRules, compactRules...))
	if err != nil {
		return nil, nil, err
	}

	wf.Rules = workflowRules
	wf.Runs = runs

	return wf, append(compactRules, runRules...), nil
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

func normalizeRules(rawRule any, currentRules []PadRule) ([]PadRule, []PadWorkflowRule, error) {
	var rule *PadRule
	var workflowRule *PadWorkflowRule
	var rules []PadRule
	var workflowRules []PadWorkflowRule

	switch r := rawRule.(type) {
	// a rule can be a plain string which can be an inline rule
	// or a name of a predefined rules
	// - '$size() < 10'
	// - small
	case string:
		rule = decodeRule(r)
		workflowRule = &PadWorkflowRule{
			Rule: rule.Name,
		}
	// we can also have a list of rules in the format
	// - '$size() < 10'
	// - small
	case []any:
		for _, ru := range r {
			processedRules, processedWorkflowRules, err := normalizeRules(ru, currentRules)
			if err != nil {
				return nil, nil, err
			}

			rules = append(rules, processedRules...)
			workflowRules = append(workflowRules, processedWorkflowRules...)
		}
	// we can also specify a rule as map with extra actions like
	// - rule: small
	//   extra-actions:
	//		- $comment("small")
	case map[string]any:
		decodedWorkflowRule, err := decodeWorkflowRule(r)
		if err != nil {
			return nil, nil, err
		}

		workflowRule = decodedWorkflowRule
		rule = decodeRule(decodedWorkflowRule.Rule)
	case nil:
		return nil, nil, nil
	default:
		return nil, nil, fmt.Errorf("unknown rule type %T", r)
	}

	if rule != nil {
		if _, exists := findRule(currentRules, rule.Name); !exists {
			rules = append(rules, *rule)
		}
	}

	if workflowRule != nil {
		workflowRules = append(workflowRules, *workflowRule)
	}

	return rules, workflowRules, nil
}

func normalizeActions(nonNormalizedActions any) ([]string, error) {
	var actions []string

	switch action := nonNormalizedActions.(type) {
	case string:
		actions = append(actions, action)
	case []any:
		for _, rawAction := range action {
			processedActions, err := normalizeActions(rawAction)
			if err != nil {
				return nil, err
			}

			actions = append(actions, processedActions...)
		}
	// we might have a workflow that doesn't have any actions
	// but only has extra actions in the workflows
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown action type: %T", action)
	}

	return actions, nil
}

func processRun(run any, currentRules []PadRule) ([]PadWorkflowRunBlock, []PadRule, error) {
	rules := []PadRule{}

	switch val := run.(type) {
	case string:
		return []PadWorkflowRunBlock{
			{
				Actions: []string{val},
			},
		}, rules, nil
	case map[string]any:
		mapBlock := PadWorkflowRunBlock{}
		if ruleBlock, ok := val["if"]; ok {
			processedRules, processedWorkflowRules, err := normalizeRules(ruleBlock, currentRules)
			if err != nil {
				return nil, nil, err
			}

			rules = append(rules, processedRules...)
			mapBlock.If = processedWorkflowRules
		}

		if thenBlock, ok := val["then"]; ok {
			thenBlocks, thenRules, err := processRun(thenBlock, append(currentRules, rules...))
			if err != nil {
				return nil, nil, err
			}

			rules = append(rules, thenRules...)
			mapBlock.Then = thenBlocks
		}

		if elseBlock, ok := val["else"]; ok {
			elseBlocks, elseRules, err := processRun(elseBlock, append(currentRules, rules...))
			if err != nil {
				return nil, nil, err
			}

			rules = append(rules, elseRules...)
			mapBlock.Else = elseBlocks
		}

		return []PadWorkflowRunBlock{
			mapBlock,
		}, rules, nil
	case []any:
		blocks := []PadWorkflowRunBlock{}
		rules := []PadRule{}
		for _, item := range val {
			switch itemVal := item.(type) {
			case string:
				// If the item is a string, add it as an action to a new block
				blocks = append(blocks, PadWorkflowRunBlock{
					Actions: []string{itemVal},
				})
			case map[string]any:
				// If the item is a map, parse it as a new PadWorkflowRunBlock
				childBlock, blockRules, err := processRun(itemVal, append(currentRules, rules...))
				if err != nil {
					return nil, nil, err
				}
				blocks = append(blocks, childBlock...)
				rules = append(rules, blockRules...)
			default:
				return nil, nil, fmt.Errorf("expected string or map, got %T", item)
			}
		}
		return blocks, rules, nil
	case nil:
		return nil, nil, nil
	default:
		return nil, nil, fmt.Errorf("unknown run type: %T", run)
	}
}
