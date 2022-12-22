// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"

	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
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

func Load(ctx context.Context, githubClient *gh.GithubClient, data []byte) (*ReviewpadFile, error) {
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

	overridenElemsByProperty := make(map[string](map[string]string))
	file, err = processExtends(ctx, githubClient, file, env, overridenElemsByProperty)
	if err != nil {
		return nil, err
	}

	for property, overridenElem := range overridenElemsByProperty {
		for elemName, triggerOverrideFile := range overridenElem {
			log.Printf("[WARN] the %s %s has been overriden after extending the configuration in %s", property, elemName, triggerOverrideFile)
		}
	}

	file, err = normalize(file, inlineRulesNormalizer())
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
	var transformedRules []PadRule
	for _, rule := range file.Rules {
		kind := rule.Kind
		if rule.Kind == "" {
			kind = "patch"
		}
		transformedRules = append(transformedRules, PadRule{
			Name:        rule.Name,
			Kind:        kind,
			Description: rule.Description,
			Spec:        transformAladinoExpression(rule.Spec),
		})
	}

	var transformedWorkflows []PadWorkflow
	for _, workflow := range file.Workflows {
		var transformedRules []PadWorkflowRule
		for _, rule := range workflow.Rules {
			var transformedExtraActions []string
			for _, extraAction := range rule.ExtraActions {
				transformedExtraActions = append(transformedExtraActions, transformAladinoExpression(extraAction))
			}

			transformedRules = append(transformedRules, PadWorkflowRule{
				Rule:         rule.Rule,
				ExtraActions: transformedExtraActions,
			})
		}

		var transformedActions []string
		for _, action := range workflow.Actions {
			transformedActions = append(transformedActions, transformAladinoExpression(action))
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
				transformedActions = append(transformedActions, transformAladinoExpression(action))
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
		Version:        file.Version,
		Edition:        file.Edition,
		Mode:           file.Mode,
		IgnoreErrors:   file.IgnoreErrors,
		MetricsOnMerge: file.MetricsOnMerge,
		Imports:        file.Imports,
		Extends:        file.Extends,
		Groups:         file.Groups,
		Rules:          transformedRules,
		Labels:         file.Labels,
		Workflows:      transformedWorkflows,
		Pipelines:      transformedPipelines,
		Recipes:        file.Recipes,
	}
}

func loadImport(reviewpadImport PadImport) (*ReviewpadFile, string, error) {
	resp, err := http.Get(reviewpadImport.Url)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
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
		file.appendPipelines(subTreeFile)
		file.appendRecipes(subTreeFile)
	}

	// reset all imports
	file.Imports = []PadImport{}

	return file, nil
}

func loadExtension(ctx context.Context, githubClient *gh.GithubClient, extension string) (*ReviewpadFile, string, error) {
	branch, filePath, err := utils.ValidateUrl(extension)
	if err != nil {
		return nil, "", err
	}

	content, err := githubClient.DownloadContents(ctx, filePath, branch)
	if err != nil {
		return nil, "", err
	}

	file, err := parse(content)
	if err != nil {
		return nil, "", err
	}

	return file, hash(content), nil
}

// processExtends inlines files into the current reviewpad file
// precedence: current file > extends file
// Post-condition: ReviewpadFile without extends statements
func processExtends(ctx context.Context, githubClient *gh.GithubClient, file *ReviewpadFile, env *LoadEnv, overridenElemsByProperty map[string](map[string]string)) (*ReviewpadFile, error) {
	extendedFile := &ReviewpadFile{}
	for _, extensionUri := range file.Extends {
		eFile, eHash, err := loadExtension(ctx, githubClient, extensionUri)
		if err != nil {
			return nil, err
		}

		checkForOverridens(overridenElemsByProperty, file, eFile, extensionUri)

		if _, ok := env.Stack[eHash]; ok {
			return nil, fmt.Errorf("loader: cyclic extends dependency")
		}

		env.Stack[eHash] = true

		extensionFile, err := processExtends(ctx, githubClient, eFile, env, overridenElemsByProperty)
		if err != nil {
			return nil, err
		}

		// remove from the stack
		delete(env.Stack, eHash)

		extendedFile.extend(extensionFile)
	}

	extendedFile.extend(file)

	// reset all extends
	extendedFile.Extends = []string{}

	return extendedFile, nil
}

func checkForOverridens(overridenElemsByProperty map[string](map[string]string), file *ReviewpadFile, triggerOverrideFile *ReviewpadFile, triggerOverrideFileUri string) {
	// Check for overriden labels
	for labelKeyName := range file.Labels {
		for eLabelKeyName := range triggerOverrideFile.Labels {
			if labelKeyName == eLabelKeyName {
				if _, found := overridenElemsByProperty["label"]; !found {
					overridenElemsByProperty["label"] = make(map[string]string)
				}
				overridenElemsByProperty["label"][eLabelKeyName] = triggerOverrideFileUri
			}
		}
	}

	// Check for overriden groups
	for _, group := range file.Groups {
		for _, eGroup := range triggerOverrideFile.Groups {
			if group.Name == eGroup.Name {
				if _, found := overridenElemsByProperty["group"]; !found {
					overridenElemsByProperty["group"] = make(map[string]string)
				}
				overridenElemsByProperty["group"][group.Name] = triggerOverrideFileUri
			}
		}
	}

	// Check for overriden rules
	for _, rule := range file.Rules {
		for _, eRule := range triggerOverrideFile.Rules {
			if rule.Name == eRule.Name {
				if _, found := overridenElemsByProperty["rule"]; !found {
					overridenElemsByProperty["rule"] = make(map[string]string)
				}
				overridenElemsByProperty["rule"][rule.Name] = triggerOverrideFileUri
			}
		}
	}

	// Check for overriden workflows
	for _, workflow := range file.Workflows {
		for _, eWorkflow := range triggerOverrideFile.Workflows {
			if workflow.Name == eWorkflow.Name {
				if _, found := overridenElemsByProperty["workflow"]; !found {
					overridenElemsByProperty["workflow"] = make(map[string]string)
				}
				overridenElemsByProperty["workflow"][workflow.Name] = triggerOverrideFileUri
			}
		}
	}

	// Check for overriden pipelines
	for _, pipeline := range file.Pipelines {
		for _, ePipeline := range triggerOverrideFile.Pipelines {
			if pipeline.Name == ePipeline.Name {
				if _, found := overridenElemsByProperty["pipeline"]; !found {
					overridenElemsByProperty["pipeline"] = make(map[string]string)
				}
				overridenElemsByProperty["pipeline"][pipeline.Name] = triggerOverrideFileUri
			}
		}
	}
}
