// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/reviewpad/reviewpad/v3/handler"
	"gopkg.in/yaml.v3"
)

type LoadEnv struct {
	Visited map[string]bool
	Stack   map[string]bool
}

func hash(data []byte) string {
	dataHash := sha256.Sum256(data)
	dHash := fmt.Sprintf("%x", dataHash)
	return dHash
}

func Load(data []byte) (*ReviewpadFile, error) {
	file, err := parse(data)
	if err != nil {
		return nil, err
	}

	dHash := hash(data)

	visited := make(map[string]bool)
	stack := make(map[string]bool)
	visited[dHash] = true
	stack[dHash] = true

	env := &LoadEnv{
		Visited: visited,
		Stack:   stack,
	}

	file, err = processImports(file, env)
	if err != nil {
		return nil, err
	}

	file, err = processInlineRules(file)
	if err != nil {
		return nil, err
	}

	return transform(file), nil
}

func parse(data []byte) (*ReviewpadFile, error) {
	file := ReviewpadFile{}
	err := yaml.Unmarshal([]byte(data), &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func transform(file *ReviewpadFile) *ReviewpadFile {
	var transformedWorkflows []PadWorkflow
	for _, workflow := range file.Workflows {
		var transformedRules []PadWorkflowRule
		for _, rule := range workflow.Rules {
			var transformedExtraActions []string
			for _, extraAction := range rule.ExtraActions {
				transformedExtraActions = append(transformedExtraActions, transformActionStr(extraAction))
			}

			transformedRules = append(transformedRules, PadWorkflowRule{
				Rule:         rule.Rule,
				ExtraActions: transformedExtraActions,
			})
		}

		var transformedActions []string
		for _, action := range workflow.Actions {
			transformedActions = append(transformedActions, transformActionStr(action))
		}

		transformedOn := []handler.TargetEntityKind{handler.PullRequest}
		if len(workflow.On) > 0 {
			transformedOn = workflow.On
		}

		transformedWorkflows = append(transformedWorkflows, PadWorkflow{
			Name:        workflow.Name,
			On:          transformedOn,
			Description: workflow.Description,
			Rules:       transformedRules,
			Actions:     transformedActions,
			AlwaysRun:   workflow.AlwaysRun,
		})
	}

	var transformedPipelines []PadPipeline

	for _, pipeline := range file.Pipelines {
		var transformedStages []PadStage

		for _, stage := range pipeline.Stages {
			var transformedActions []string
			for _, action := range stage.Actions {
				transformedActions = append(transformedActions, transformActionStr(action))
			}

			transformedStages = append(transformedStages, PadStage{
				Actions: transformedActions,
				Until:   stage.Until,
			})
		}

		transformedPipelines = append(transformedPipelines, PadPipeline{
			Name:        pipeline.Name,
			Description: pipeline.Description,
			Trigger:     pipeline.Trigger,
			Stages:      transformedStages,
		})
	}

	return &ReviewpadFile{
		Version:      file.Version,
		Edition:      file.Edition,
		Mode:         file.Mode,
		IgnoreErrors: file.IgnoreErrors,
		Imports:      file.Imports,
		Groups:       file.Groups,
		Rules:        file.Rules,
		Labels:       file.Labels,
		Workflows:    transformedWorkflows,
		Pipelines:    transformedPipelines,
	}
}

func loadImport(reviewpadImport PadImport) (*ReviewpadFile, string, error) {
	resp, err := http.Get(reviewpadImport.Url)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	file, err := parse(content)
	if err != nil {
		return nil, "", err
	}

	return file, hash(content), nil
}

// processImports inlines the imports files into the current reviewpad file
// Post-condition: ReviewpadFile without import statements
func processImports(file *ReviewpadFile, env *LoadEnv) (*ReviewpadFile, error) {
	for _, reviewpadImport := range file.Imports {
		iFile, idHash, err := loadImport(reviewpadImport)
		if err != nil {
			return nil, err
		}

		// check for cycles
		if _, ok := env.Stack[idHash]; ok {
			return nil, fmt.Errorf("loader: cyclic dependency")
		}

		// optimize visits
		if _, ok := env.Visited[idHash]; ok {
			continue
		}

		// DFS call inline imports
		// update the environment
		env.Stack[idHash] = true
		env.Visited[idHash] = true

		subTreeFile, err := processImports(iFile, env)
		if err != nil {
			return nil, err
		}

		// remove from the stack
		delete(env.Stack, idHash)

		// append labels, rules and workflows
		file.appendLabels(subTreeFile)
		file.appendGroups(subTreeFile)
		file.appendRules(subTreeFile)
		file.appendWorkflows(subTreeFile)
	}

	// reset all imports
	file.Imports = []PadImport{}

	return file, nil
}

// processInlineRules normalizes the inline rules in the file
// by converting the inline rules into a PadWorkflowRule
func processInlineRules(file *ReviewpadFile) (*ReviewpadFile, error) {
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
	alreadyGeneratedRules := map[string]bool{}

	for i, workflow := range reviewpadFile.Workflows {
		processedWorkflow, rules, err := processInlineRulesOnWorkflow(workflow, file.Rules, alreadyGeneratedRules)
		if err != nil {
			return nil, err
		}

		reviewpadFile.Rules = append(reviewpadFile.Rules, rules...)
		reviewpadFile.Workflows[i] = *processedWorkflow
	}

	return reviewpadFile, nil
}

func processInlineRulesOnWorkflow(workflow PadWorkflow, currentRules []PadRule, alreadyGeneratedRules map[string]bool) (*PadWorkflow, []PadRule, error) {
	wf := &PadWorkflow{
		Name:        workflow.Name,
		Description: workflow.Description,
		AlwaysRun:   workflow.AlwaysRun,
		Rules:       workflow.Rules,
		Actions:     workflow.Actions,
		On:          workflow.On,
	}
	currentNormalizedRules := make([]PadRule, 0)

	for _, rule := range workflow.NonNormalizedRules {
		switch r := rule.(type) {
		case string:
			padRule := inlineRuleToPadRule(r)

			if !alreadyGeneratedRules[padRule.Name] {
				alreadyGeneratedRules[padRule.Name] = true
				currentNormalizedRules = append(currentNormalizedRules, padRule)
			}

			wf.Rules = append(wf.Rules, PadWorkflowRule{
				Rule: padRule.Name,
			})
		case map[string]interface{}:
			rule, err := mapToPadWorkflowRule(r)
			if err != nil {
				return nil, nil, err
			}

			// if the rule doesn't exist in the list of available rules
			// treat it as an inline rule with extra actions
			if _, exists := findRule(currentRules, rule.Rule); !exists {
				padRule := inlineRuleToPadRule(rule.Rule)

				if !alreadyGeneratedRules[padRule.Name] {
					alreadyGeneratedRules[padRule.Name] = true
					currentNormalizedRules = append(currentNormalizedRules, padRule)
				}

				wf.Rules = append(wf.Rules, PadWorkflowRule{
					Rule:         padRule.Name,
					ExtraActions: rule.ExtraActions,
				})
			} else {
				wf.Rules = append(wf.Rules, *rule)
			}
		default:
			return nil, nil, fmt.Errorf("unknown rule type %T", r)
		}
	}

	// we are removing the non-normalized rules from the workflow
	// to make things a bit easier to test
	workflow.NonNormalizedRules = nil

	return wf, currentNormalizedRules, nil
}

func mapToPadWorkflowRule(rule map[string]interface{}) (*PadWorkflowRule, error) {
	r := &PadWorkflowRule{}
	err := mapstructure.Decode(rule, r)
	return r, err
}

func inlineRuleToPadRule(rule string) PadRule {
	return PadRule{
		Name: fmt.Sprintf("inline rule %s", rule),
		Spec: rule,
		Kind: "patch",
	}
}
