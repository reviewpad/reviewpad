// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/sirupsen/logrus"
)

func getAllMatches(pattern string, groups []PadGroup, rules []PadRule, workflows []PadWorkflow) []string {
	rePatternFnCall := regexp.MustCompile(pattern)
	allGroupFunctionCalls := make([]string, 0)
	for _, group := range groups {
		spec := group.Spec
		groupFunctionCalls := rePatternFnCall.FindAllString(spec, -1)

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	for _, rule := range rules {
		spec := rule.Spec
		groupFunctionCalls := rePatternFnCall.FindAllString(spec, -1)

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	for _, workflow := range workflows {
		actions := workflow.Actions
		groupFunctionCalls := make([]string, 0)
		for _, action := range actions {
			groupFunctionCalls = append(groupFunctionCalls, rePatternFnCall.FindAllString(action, -1)...)
		}

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	return allGroupFunctionCalls
}

// Validations:
// - Every rule has a (unique) name
// - Every rule has a kind
// - Every rules has a spec
func lintRules(padRules []PadRule) error {
	rulesName := make([]string, 0)

	for _, rule := range padRules {
		if rule.Name == "" {
			return fmt.Errorf("rule %v has invalid name", rule)
		}

		for _, ruleName := range rulesName {
			if ruleName == rule.Name {
				return fmt.Errorf("rule with the name %v already exists", rule.Name)
			}
		}

		ruleKind := rule.Kind
		if !utils.ElementOf(kinds, ruleKind) {
			return fmt.Errorf("rule %v has invalid kind %v", rule.Name, ruleKind)
		}

		if rule.Spec == "" {
			return fmt.Errorf("rule %v has empty spec", rule.Name)
		}

		rulesName = append(rulesName, rule.Name)
	}

	return nil
}

// Validations:
// - Group has unique name
func lintGroups(log *logrus.Entry, padGroups []PadGroup) error {
	groupsName := make([]string, 0)

	for _, group := range padGroups {
		log.Infof("analyzing group %v", group.Name)

		if group.Name == "" {
			return fmt.Errorf("group %v has invalid name", group)
		}

		for _, groupName := range groupsName {
			if groupName == group.Name {
				return fmt.Errorf("group with the name %v already exists", group.Name)
			}
		}

		groupsName = append(groupsName, group.Name)
	}

	return nil
}

// Validations:
// - Workflow has unique name
// - Workflow has non empty rules
// - Workflow has only known rules
// - Workflow run block has a then block or actions
func lintWorkflows(log *logrus.Entry, rules []PadRule, padWorkflows []PadWorkflow) error {
	workflowsName := make([]string, 0)
	workflowHasExtraActions := false

	for _, workflow := range padWorkflows {
		log.Infof("linting workflow `%v`", workflow.Name)

		workflowHasActions := len(workflow.Actions) > 0

		for _, workflowName := range workflowsName {
			if workflowName == workflow.Name {
				return fmt.Errorf("workflow with the name `%v` already exists", workflow.Name)
			}
		}

		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			if ruleName == "" {
				return fmt.Errorf("workflow has an empty rule")
			}

			_, exists := findRule(rules, ruleName)
			if !exists {
				return fmt.Errorf("rule `%v` is unknown", ruleName)
			}

			workflowHasExtraActions = len(rule.ExtraActions) > 0
			if !workflowHasExtraActions && !workflowHasActions {
				log.Warnf("rule `%v` will be ignored since it has no actions", ruleName)
			}
		}

		if !workflowHasActions && !workflowHasExtraActions {
			log.Debug("workflow has no actions")
		}

		workflowsName = append(workflowsName, workflow.Name)

		for _, run := range workflow.Runs {
			if err := validateWorkflowRun(&run, &workflow); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateWorkflowRun(run *PadWorkflowRunBlock, workflow *PadWorkflow) error {
	var hasForEachBlock = run.ForEach != nil
	var hasActions = run.Actions != nil && len(run.Actions) > 0
	var hasThenActions = run.Then != nil && len(run.Then) > 0
	var hasExtraActions bool
	for _, rule := range run.If {
		hasExtraActions = len(rule.ExtraActions) > 0
		if hasExtraActions {
			break
		}
	}

	// Old style workflow if/then/else are converted to new run block format.
	// The old style workflow allows if blocks to have extra actions.
	// Because of this, a run block can have extra actions.
	if !hasThenActions && !hasActions && !hasExtraActions && !hasForEachBlock {
		return fmt.Errorf("workflow `%v` has a run block without a 'then' block, no actions, no extra actions or a for each block", workflow.Name)
	}

	for _, thenRun := range run.Then {
		err := validateWorkflowRun(&thenRun, workflow)
		if err != nil {
			return err
		}
	}

	for _, elseRun := range run.Else {
		err := validateWorkflowRun(&elseRun, workflow)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validations
// - Check that all rules are being used
// - Check that all referenced rules exist
func lintRulesMentions(log *logrus.Entry, rules []PadRule, groups []PadGroup, workflows []PadWorkflow) error {
	totalUsesByRule := make(map[string]int, len(rules))

	for _, rule := range rules {
		totalUsesByRule[rule.Name] = 0
	}

	for _, workflow := range workflows {
		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			_, exists := findRule(rules, ruleName)

			if exists {
				totalUsesByRule[ruleName]++
			}
		}
	}

	for _, ruleName := range getCallsToRuleBuiltIn(groups, rules, workflows) {
		_, ok := findRule(rules, ruleName)
		if !ok {
			return fmt.Errorf("the rule `%v` isn't defined", ruleName)
		}
		totalUsesByRule[ruleName]++
	}

	for ruleName, totalUses := range totalUsesByRule {
		if totalUses == 0 {
			log.Warnf("unused rule `%v`", ruleName)
		}
	}

	return nil
}

func getCallsToRuleBuiltIn(groups []PadGroup, rules []PadRule, workflows []PadWorkflow) []string {
	allRuleFunctionCalls := make([]string, 0)

	gotFunctionCalls := getAllMatches(`\$rule\("[^)]*"\)`, groups, rules, workflows)

	reRuleMention := regexp.MustCompile(`"(.*?)"`)
	for _, ruleWithRuleCall := range gotFunctionCalls {
		for _, ruleCall := range reRuleMention.FindAllString(ruleWithRuleCall, -1) {
			ruleName := ruleCall[1 : len(ruleCall)-1]

			allRuleFunctionCalls = append(allRuleFunctionCalls, ruleName)
		}
	}

	return allRuleFunctionCalls
}

// Validations
// - Check that all groups are being used
// - Check that all referenced groups exist
func lintGroupsMentions(groups []PadGroup, rules []PadRule, workflows []PadWorkflow) error {
	allGroupFunctionCalls := getAllMatches(`\$group\(".*"\)`, groups, rules, workflows)

	reGroupMention := regexp.MustCompile(`"(.*?)"`)
	for _, groupFunctionCall := range allGroupFunctionCalls {
		groupMention := reGroupMention.FindString(groupFunctionCall)
		// Remove quotation marks
		groupMention = groupMention[1 : len(groupMention)-1]

		_, ok := findGroup(groups, groupMention)
		if !ok {
			return fmt.Errorf("the group `%v` isn't defined", groupMention)
		}
	}

	return nil
}

func lintShadowedVariablesInRuns(runs []PadWorkflowRunBlock, definedVariables map[string]bool) error {
	for _, run := range runs {
		if run.ForEach != nil {
			if _, ok := definedVariables[run.ForEach.Value]; ok {
				return fmt.Errorf("variable shadowing is not allowed: the variable `%s` is already defined", run.ForEach.Value)
			}

			definedVariables[run.ForEach.Value] = true

			err := lintShadowedVariablesInRuns(run.ForEach.Do, definedVariables)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func lintShadowedVariables(workflows []PadWorkflow) error {
	for _, workflow := range workflows {
		definedVariables := map[string]bool{}
		err := lintShadowedVariablesInRuns(workflow.Runs, definedVariables)
		if err != nil {
			return err
		}
	}

	return nil
}

func Lint(file *ReviewpadFile, logger *logrus.Entry) error {
	err := lintGroups(logger, file.Groups)
	if err != nil {
		return err
	}

	err = lintRules(file.Rules)
	if err != nil {
		return err
	}

	err = lintWorkflows(logger, file.Rules, file.Workflows)
	if err != nil {
		return err
	}

	err = lintRulesMentions(logger, file.Rules, file.Groups, file.Workflows)
	if err != nil {
		return err
	}

	err = lintShadowedVariables(file.Workflows)
	if err != nil {
		return err
	}

	return lintGroupsMentions(file.Groups, file.Rules, file.Workflows)
}
