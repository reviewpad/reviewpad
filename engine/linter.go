// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"

	"github.com/reviewpad/reviewpad/utils"
	"github.com/reviewpad/reviewpad/utils/fmtio"
)

func lintError(format string, a ...interface{}) error {
	return fmtio.Errorf("lint", format, a...)
}

func lintLog(format string, a ...interface{}) {
	fmtio.LogPrintln("lint", format, a...)
}

// Validations:
// - Every rule has a kind
// - Every rules has a spec
func lintRules(padRules map[string]PadRule) error {
	for ruleName, rule := range padRules {
		ruleKind := rule.Kind
		if !utils.ElementOf(kinds, ruleKind) {
			return lintError("rule %v has invalid kind %v", ruleName, ruleKind)
		}

		if rule.Spec == "" {
			return lintError("rule %v has empty spec", ruleName)
		}
	}

	return nil
}

// Validations:
// - Workflow has unique name
// - Workflow has rules
// - Workflow has non empty rules
// - Workflow has only known rules
func lintWorkflows(rules map[string]PadRule, padWorkflows []PadWorkflow) error {
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

			_, exists := rules[ruleName]
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
// - Check that all rules are being used on workflows
func lintRulesWithWorkflows(rules map[string]PadRule, padWorkflows []PadWorkflow) error {
	totalUsesByRule := make(map[string]int, len(rules))

	for ruleName := range rules {
		totalUsesByRule[ruleName] = 0
	}

	for _, workflow := range padWorkflows {
		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			_, exists := rules[ruleName]

			if exists {
				totalUsesByRule[ruleName]++
			}
		}
	}

	for ruleName, totalUses := range totalUsesByRule {
		if totalUses == 0 {
			return lintError("unused rule %v", ruleName)
		}
	}

	return nil
}

func lintGroupsMentions(groups map[string]PadGroup, rules map[string]PadRule, workflows []PadWorkflow) error {
	reGroupFnCall := regexp.MustCompile(`\$group\(".*"\)`)
	allGroupFunctionCalls := make([]string, 0)
	for _, group := range groups {
		spec := group.Spec
		groupFunctionCalls := reGroupFnCall.FindAllString(spec, -1)

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	for _, rule := range rules {
		spec := rule.Spec
		groupFunctionCalls := reGroupFnCall.FindAllString(spec, -1)

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	for _, workflow := range workflows {
		actions := workflow.Actions
		groupFunctionCalls := make([]string, 0)
		for _, action := range actions {
			groupFunctionCalls = append(groupFunctionCalls, reGroupFnCall.FindAllString(action, -1)...)
		}

		if len(groupFunctionCalls) == 0 {
			continue
		}

		allGroupFunctionCalls = append(allGroupFunctionCalls, groupFunctionCalls...)
	}

	reGroupMention := regexp.MustCompile(`"(.*?)"`)
	for _, groupFunctionCall := range allGroupFunctionCalls {
		groupMention := reGroupMention.FindString(groupFunctionCall)
		// Remove quotation marks
		groupMention = groupMention[1 : len(groupMention)-1]

		_, ok := groups[groupMention]
		if !ok {
			return lintError("the group %v isn't defined", groupMention)
		}
	}

	return nil
}

func Lint(file *ReviewpadFile) error {
	err := lintRules(file.Rules)
	if err != nil {
		return err
	}

	err = lintWorkflows(file.Rules, file.Workflows)
	if err != nil {
		return err
	}

	// skipping this validation as it is now possible for rules to reference each other
	// TODO: #8 Improve this validation
	// err = lintRulesWithWorkflows(file.Rules, file.Workflows)
	// if err != nil {
	// 	return err
	// }

	return lintGroupsMentions(file.Groups, file.Rules, file.Workflows)
}
