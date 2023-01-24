// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type LoadEnv struct {
	Visited map[string]bool
	Stack   map[string]bool
}

func logExtendedProperties(logger *logrus.Entry, extensionUri string, extensionFile, previousResultFile, resultFile *ReviewpadFile) {
	// Extended labels
	for labelKeyName, label := range resultFile.Labels {
		if _, ok := extensionFile.Labels[labelKeyName]; ok {
			if _, ok := previousResultFile.Labels[labelKeyName]; ok {
				if extensionFile.Labels[labelKeyName].equals(label) {
					logger.Warnf("label %s has been overridden by %s", labelKeyName, extensionUri)
				}
			}
		}
	}

	// Extended groups
	for _, group := range resultFile.Groups {
		var extensionFileGroup, previousResultFileGroup PadGroup
		for _, groupA := range extensionFile.Groups {
			if groupA.Name == group.Name {
				extensionFileGroup = groupA
				break
			}
		}

		for _, groupB := range previousResultFile.Groups {
			if groupB.Name == group.Name {
				previousResultFileGroup = groupB
				break
			}
		}

		if !extensionFileGroup.equals(PadGroup{}) && !previousResultFileGroup.equals(PadGroup{}) {
			if extensionFileGroup.equals(group) {
				logger.Warnf("group %s has been overridden by %s", group.Name, extensionUri)
			}
		}
	}

	// Extended rules
	for _, rule := range resultFile.Rules {
		var extensionFileRule, previousResultFileRule PadRule
		for _, ruleA := range extensionFile.Rules {
			if ruleA.Name == rule.Name {
				extensionFileRule = ruleA
				break
			}
		}

		for _, ruleB := range previousResultFile.Rules {
			if ruleB.Name == rule.Name {
				previousResultFileRule = ruleB
				break
			}
		}

		if !extensionFileRule.equals(PadRule{}) && !previousResultFileRule.equals(PadRule{}) {
			if extensionFileRule.equals(rule) {
				logger.Warnf("rule %s has been overridden by %s", rule.Name, extensionUri)
			}
		}
	}

	// Extended workflows
	for _, workflow := range resultFile.Workflows {
		var extensionFileWorkflow, previousResultFileWorkflow PadWorkflow
		for _, workflowA := range extensionFile.Workflows {
			if workflowA.Name == workflow.Name {
				extensionFileWorkflow = workflowA
				break
			}
		}

		for _, workflowB := range previousResultFile.Workflows {
			if workflowB.Name == workflow.Name {
				previousResultFileWorkflow = workflowB
				break
			}
		}

		if !extensionFileWorkflow.equals(PadWorkflow{}) && !previousResultFileWorkflow.equals(PadWorkflow{}) {
			if extensionFileWorkflow.equals(workflow) {
				logger.Warnf("workflow %s has been overridden by %s", workflow.Name, extensionUri)
			}
		}
	}

	// Extended pipelines
	for _, pipeline := range resultFile.Pipelines {
		var extensionFilePipeline, previousResultFilePipeline PadPipeline
		for _, pipelineA := range extensionFile.Pipelines {
			if pipelineA.Name == pipeline.Name {
				extensionFilePipeline = pipelineA
				break
			}
		}

		for _, pipelineB := range previousResultFile.Pipelines {
			if pipelineB.Name == pipeline.Name {
				previousResultFilePipeline = pipelineB
				break
			}
		}

		if !extensionFilePipeline.equals(PadPipeline{}) && !previousResultFilePipeline.equals(PadPipeline{}) {
			if extensionFilePipeline.equals(pipeline) {
				logger.Warnf("pipeline %s has been overridden by %s", pipeline.Name, extensionUri)
			}
		}
	}
}

func hash(data []byte) string {
	dataHash := sha256.Sum256(data)
	dHash := fmt.Sprintf("%x", dataHash)
	return dHash
}

func Load(ctx context.Context, log *logrus.Entry, githubClient *gh.GithubClient, data []byte) (*ReviewpadFile, error) {
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

	file, err = processExtends(ctx, log, githubClient, file, "root file", env)
	if err != nil {
		return nil, err
	}

	log.Infof("FILE: %+v", file)

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
func processExtends(ctx context.Context, log *logrus.Entry, githubClient *gh.GithubClient, file *ReviewpadFile, fileUri string, env *LoadEnv) (*ReviewpadFile, error) {
	extendedFile := &ReviewpadFile{}
	for _, extensionUri := range file.Extends {
		eFile, eHash, err := loadExtension(ctx, githubClient, extensionUri)
		if err != nil {
			return nil, err
		}

		if _, ok := env.Stack[eHash]; ok {
			return nil, fmt.Errorf("loader: cyclic extends dependency")
		}

		env.Stack[eHash] = true

		extensionFile, err := processExtends(ctx, log, githubClient, eFile, extensionUri, env)
		if err != nil {
			return nil, err
		}

		// remove from the stack
		delete(env.Stack, eHash)

		previousExtendedFile := *extendedFile

		extendedFile.extend(extensionFile)

		logExtendedProperties(log, extensionUri, extensionFile, &previousExtendedFile, extendedFile)
	}

	previousExtendedFile := *extendedFile

	extendedFile.extend(file)

	logExtendedProperties(log, fileUri, file, &previousExtendedFile, extendedFile)

	// reset all extends
	extendedFile.Extends = []string{}

	return extendedFile, nil
}
