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

	"github.com/google/go-github/v49/github"
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

func logExtendedProperties(logger *logrus.Entry, fileA *ReviewpadFile, fileAUri string, fileB, resultFile *ReviewpadFile) {
	// Extended labels
	for labelName, label := range resultFile.Labels {
		if _, ok := fileA.Labels[labelName]; ok {
			if _, ok := fileB.Labels[labelName]; ok {
				if fileA.Labels[labelName].equals(label) {
					logger.Warnf("label '%s' has been overridden by %s", labelName, fileAUri)
				}
			}
		}
	}

	// Extended groups
	for _, group := range resultFile.Groups {
		var groupA, groupB *PadGroup
		for _, gA := range fileA.Groups {
			if gA.Name == group.Name {
				groupA = &gA
				break
			}
		}

		for _, gB := range fileB.Groups {
			if gB.Name == group.Name {
				groupB = &gB
				break
			}
		}

		if groupA != nil && groupB != nil {
			if groupA.equals(group) {
				logger.Warnf("group '%s' has been overridden by %s", group.Name, fileAUri)
			}
		}
	}

	// Extended rules
	for _, rule := range resultFile.Rules {
		var ruleA, ruleB *PadRule
		for _, rA := range fileA.Rules {
			if rA.Name == rule.Name {
				ruleA = &rA
				break
			}
		}

		for _, rB := range fileB.Rules {
			if rB.Name == rule.Name {
				ruleB = &rB
				break
			}
		}

		if ruleA != nil && ruleB != nil {
			if ruleA.equals(rule) {
				logger.Warnf("rule '%s' has been overridden by %s", rule.Name, fileAUri)
			}
		}
	}

	// Extended workflows
	for _, workflow := range resultFile.Workflows {
		var workflowA, workflowB *PadWorkflow
		for _, wA := range fileA.Workflows {
			if wA.Name == workflow.Name {
				workflowA = &wA
				break
			}
		}

		for _, wB := range fileB.Workflows {
			if wB.Name == workflow.Name {
				workflowB = &wB
				break
			}
		}

		if workflowA != nil && workflowB != nil {
			if workflowA.equals(workflow) {
				logger.Warnf("workflow '%s' has been overridden by %s", workflow.Name, fileAUri)
			}
		}
	}

	// Extended pipelines
	for _, pipeline := range resultFile.Pipelines {
		var pipelineA, pipelineB *PadPipeline
		for _, pA := range fileA.Pipelines {
			if pA.Name == pipeline.Name {
				pipelineA = &pA
				break
			}
		}

		for _, pB := range fileB.Pipelines {
			if pB.Name == pipeline.Name {
				pipelineB = &pB
				break
			}
		}

		if pipelineA != nil && pipelineB != nil {
			if pipelineA.equals(pipeline) {
				logger.Warnf("pipeline '%s' has been overridden by %s", pipeline.Name, fileAUri)
			}
		}
	}
}

func hash(data []byte) string {
	dataHash := sha256.Sum256(data)
	dHash := fmt.Sprintf("%x", dataHash)
	return dHash
}

func Load(ctx context.Context, logger *logrus.Entry, githubClient *gh.GithubClient, data []byte) (*ReviewpadFile, error) {
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

	file, err = processExtends(ctx, logger, githubClient, file, "root file", env)
	if err != nil {
		return nil, err
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
		if responseErr, ok := err.(*github.ErrorResponse); ok && responseErr.Response.StatusCode == 404 {
			contactUs := "If you still encounter this error, please reach out to us at #help on https://reviewpad.com/discord for additional support."
			switch responseErr.Message {
			case "Repository access blocked":
				return nil, "", fmt.Errorf("we encountered an authorization error while processing the 'extends' property in your Reviewpad configuration. Please ensure that you have the Reviewpad GitHub App installed in all repositories you are trying to extend from. %s", contactUs)
			default:
				return nil, "", fmt.Errorf("we encountered an error while processing the 'extends' property in your Reviewpad configuration. This problem may be due to an incorrect URL or unauthenticated access. Please ensure that you have Reviewpad GitHub App installed in all repositories you are trying to extend from and that the file URL is correct. %s", contactUs)
			}
		}
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
func processExtends(ctx context.Context, logger *logrus.Entry, githubClient *gh.GithubClient, file *ReviewpadFile, fileUri string, env *LoadEnv) (*ReviewpadFile, error) {
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

		extensionFile, err := processExtends(ctx, logger, githubClient, eFile, extensionUri, env)
		if err != nil {
			return nil, err
		}

		// remove from the stack
		delete(env.Stack, eHash)

		extendedFileBeforeExtend := *extendedFile

		extendedFile.extend(extensionFile)

		logExtendedProperties(logger, extensionFile, extensionUri, &extendedFileBeforeExtend, extendedFile)
	}

	extendedFileBeforeExtend := *extendedFile

	extendedFile.extend(file)

	logExtendedProperties(logger, file, fileUri, &extendedFileBeforeExtend, extendedFile)

	// reset all extends
	extendedFile.Extends = []string{}

	return extendedFile, nil
}
