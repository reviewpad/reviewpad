// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"

	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/reviewpad/reviewpad/v2/utils/fmtio"
)

func lintError(format string, a ...interface{}) error {
	return fmtio.Errorf("lint", format, a...)
}

func lintLog(format string, a ...interface{}) {
	fmtio.LogPrintln("lint", format, a...)
}

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
			return lintError("rule %v has invalid name", rule)
		}

		for _, ruleName := range rulesName {
			if ruleName == rule.Name {
				return lintError("rule with the name %v already exists", rule.Name)
			}
		}

		ruleKind := rule.Kind
		if !utils.ElementOf(kinds, ruleKind) {
			return lintError("rule %v has invalid kind %v", rule.Name, ruleKind)
		}

		if rule.Spec == "" {
			return lintError("rule %v has empty spec", rule.Name)
		}

		rulesName = append(rulesName, rule.Name)
	}

	return nil
}

// Validations:
// - Group has unique name
func lintGroups(padGroups []PadGroup) error {
	groupsName := make([]string, 0)

	for _, group := range padGroups {
		lintLog("analyzing group %v", group.Name)

		if group.Name == "" {
			return lintError("group %v has invalid name", group)
		}

		for _, groupName := range groupsName {
			if groupName == group.Name {
				return lintError("group with the name %v already exists", group.Name)
			}
		}

		groupsName = append(groupsName, group.Name)
	}

	return nil
}

// Validations:
// - Workflow has unique name
// - Workflow has rules
// - Workflow has non empty rules
// - Workflow has only known rules
func lintWorkflows(rules []PadRule, padWorkflows []PadWorkflow) error {
	workflowsName := make([]string, 0)
	workflowHasExtraActions := false

	for _, workflow := range padWorkflows {
		lintLog("analyzing workflow %v", workflow.Name)

		workflowHasActions := len(workflow.Actions) > 0

		for _, workflowName := range workflowsName {
			if workflowName == workflow.Name {
				return lintError("workflow with the name %v already exists", workflow.Name)
			}
		}

		if len(workflow.Rules) == 0 {
			return lintError("workflow %v does not have rules", workflow.Name)
		}

		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			if ruleName == "" {
				return lintError("workflow has an empty rule")
			}

			_, exists := findRule(rules, ruleName)
			if !exists {
				return lintError("rule %v is unknown", ruleName)
			}

			workflowHasExtraActions = len(rule.ExtraActions) > 0
			if !workflowHasExtraActions && !workflowHasActions {
				lintLog("warning: rule %v will be ignored since it has no actions", ruleName)
			}
		}

		if !workflowHasActions && !workflowHasExtraActions {
			lintLog("warning: workflow has no actions")
		}

		workflowsName = append(workflowsName, workflow.Name)
	}

	return nil
}

// Validations
// - Check that all rules are being used
// - Check that all referenced rules exist
func lintRulesMentions(rules []PadRule, groups []PadGroup, workflows []PadWorkflow) error {
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

	allRuleFunctionCalls := getAllMatches(`\$rule\(".*"\)`, groups, rules, workflows)
	reRuleMention := regexp.MustCompile(`"(.*?)"`)
	for _, ruleFunctionCall := range allRuleFunctionCalls {
		ruleName := reRuleMention.FindString(ruleFunctionCall)
		// Remove quotation marks
		ruleName = ruleName[1 : len(ruleName)-1]

		_, ok := findRule(rules, ruleName)
		if !ok {
			return lintError("the rule %v isn't defined", ruleName)
		}
		totalUsesByRule[ruleName]++
	}

	for ruleName, totalUses := range totalUsesByRule {
		if totalUses == 0 {
			return lintError("unused rule %v", ruleName)
		}
	}

	return nil
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
			return lintError("the group %v isn't defined", groupMention)
		}
	}

	return nil
}

func Lint(file *ReviewpadFile) error {
	err := lintGroups(file.Groups)
	if err != nil {
		return err
	}

	err = lintRules(file.Rules)
	if err != nil {
		return err
	}

	err = lintWorkflows(file.Rules, file.Workflows)
	if err != nil {
		return err
	}

	err = lintRulesMentions(file.Rules, file.Groups, file.Workflows)
	if err != nil {
		return err
	}

	return lintGroupsMentions(file.Groups, file.Rules, file.Workflows)
}
