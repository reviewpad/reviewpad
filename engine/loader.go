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

type ReviewpadFileWithUri struct {
	uri  string
	file *ReviewpadFile
}

func logOverrides(log *logrus.Entry, extensionFilesStack []ReviewpadFileWithUri) {
	if len(extensionFilesStack) == 1 {
		return
	}

	for i := len(extensionFilesStack) - 1; i >= 1; i-- {
		extendedUri := extensionFilesStack[i].uri
		extendedFile := extensionFilesStack[i].file

		extenderUri := extensionFilesStack[i-1].uri
		if extenderUri == "" {
			extenderUri = "the root configuration"
		}
		extenderFile := extensionFilesStack[i-1].file

		overridenLabels := getOverridenLabels(extenderFile, extendedFile)
		for _, overridenLabel := range overridenLabels {
			log.Warnf("the label %s in %s has been overriden by the label in %s", overridenLabel, extendedUri, extenderUri)
		}

		overridenGroups := getOverridenGroups(extenderFile, extendedFile)
		for _, overridenGroup := range overridenGroups {
			log.Warnf("the group %s in %s has been overriden by the group in %s", overridenGroup, extendedUri, extenderUri)
		}

		overridenRules := getOverridenRules(extenderFile, extendedFile)
		for _, overridenRule := range overridenRules {
			log.Warnf("the rule %s in %s has been overriden by the rule in %s", overridenRule, extendedUri, extenderUri)
		}

		overridenWorkflows := getOverridenWorkflows(extenderFile, extendedFile)
		for _, overridenWorkflow := range overridenWorkflows {
			log.Warnf("the workflow %s in %s has been overriden by the workflow in %s", overridenWorkflow, extendedUri, extenderUri)
		}

		overridenPipelines := getOverridenPipelines(extenderFile, extendedFile)
		for _, overridenPipeline := range overridenPipelines {
			log.Warnf("the pipeline %s in %s has been overriden by the pipeline in %s", overridenPipeline, extendedUri, extenderUri)
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

	extensionFilesStack := []ReviewpadFileWithUri{
		{
			uri:  "",
			file: file,
		},
	}
	file, err = processExtends(ctx, log, githubClient, file, env, &extensionFilesStack)
	if err != nil {
		return nil, err
	}

	logOverrides(log, extensionFilesStack)

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
func processExtends(ctx context.Context, logger *logrus.Entry, githubClient *gh.GithubClient, file *ReviewpadFile, env *LoadEnv, extendedFiles *[]ReviewpadFileWithUri) (*ReviewpadFile, error) {
	extendedFile := &ReviewpadFile{}
	for _, extensionUri := range file.Extends {
		eFile, eHash, err := loadExtension(ctx, githubClient, extensionUri)
		if err != nil {
			return nil, err
		}

		*extendedFiles = append(*extendedFiles, ReviewpadFileWithUri{
			uri:  extensionUri,
			file: eFile,
		})

		if _, ok := env.Stack[eHash]; ok {
			return nil, fmt.Errorf("loader: cyclic extends dependency")
		}

		env.Stack[eHash] = true

		extensionFile, err := processExtends(ctx, logger, githubClient, eFile, env, extendedFiles)
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

func getOverridenLabels(extenderFile *ReviewpadFile, extendedFile *ReviewpadFile) []string {
	overridenLabels := make([]string, 0)
	for extendedFileLabelName := range extendedFile.Labels {
		for extenderFileLabelKeyName := range extenderFile.Labels {
			if extendedFileLabelName == extenderFileLabelKeyName {
				overridenLabels = append(overridenLabels, extendedFileLabelName)
			}
		}
	}

	return overridenLabels
}

func getOverridenGroups(extenderFile *ReviewpadFile, extendedFile *ReviewpadFile) []string {
	overridenGroups := make([]string, 0)
	for _, extendedFileGroup := range extendedFile.Groups {
		for _, extenderFileGroup := range extenderFile.Groups {
			if extendedFileGroup.Name == extenderFileGroup.Name {
				overridenGroups = append(overridenGroups, extendedFileGroup.Name)
			}
		}
	}

	return overridenGroups
}

func getOverridenRules(extenderFile *ReviewpadFile, extendedFile *ReviewpadFile) []string {
	overridenRules := make([]string, 0)
	for _, extendedFileRule := range extendedFile.Rules {
		for _, extenderFileRule := range extenderFile.Rules {
			if extendedFileRule.Name == extenderFileRule.Name {
				overridenRules = append(overridenRules, extendedFileRule.Name)
			}
		}
	}

	return overridenRules
}

func getOverridenWorkflows(extenderFile *ReviewpadFile, extendedFile *ReviewpadFile) []string {
	overridenWorkflows := make([]string, 0)
	for _, extendedFileWorkflow := range extendedFile.Workflows {
		for _, extenderFileWorkflow := range extenderFile.Workflows {
			if extendedFileWorkflow.Name == extenderFileWorkflow.Name {
				overridenWorkflows = append(overridenWorkflows, extendedFileWorkflow.Name)
			}
		}
	}

	return overridenWorkflows
}

func getOverridenPipelines(extenderFile *ReviewpadFile, extendedFile *ReviewpadFile) []string {
	overridenPipelines := make([]string, 0)
	for _, extendedFilePipeline := range extendedFile.Pipelines {
		for _, extenderFilePipeline := range extenderFile.Pipelines {
			if extendedFilePipeline.Name == extenderFilePipeline.Name {
				overridenPipelines = append(overridenPipelines, extendedFilePipeline.Name)
			}
		}
	}

	return overridenPipelines
}
