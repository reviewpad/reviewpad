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
// - Protection gate has unique name
// - Protection gate has rules
// - Protection gate has non empty rules
// - Protetion gate has only known rules
func lintGates(rules map[string]PadRule, padGates []PadProtectionGate) error {
	protectionGatesName := make([]string, 0)
	gateHasExtraActions := false

	for _, gate := range padGates {
		lintLog("analyzing protetion gate %v", gate.Name)

		gateHasActions := len(gate.Actions) > 0

		for _, gateName := range protectionGatesName {
			if gateName == gate.Name {
				return lintError("gate with the name %v already exists", gate.Name)
			}
		}

		if len(gate.PatchRules) == 0 {
			return lintError("gate %v does not have rules", gate.Name)
		}

		for _, rule := range gate.PatchRules {
			ruleName := rule.Rule
			if ruleName == "" {
				return lintError("gate has an empty rule")
			}

			_, exists := rules[ruleName]
			if !exists {
				return lintError("rule %v is unknown", ruleName)
			}

			gateHasExtraActions = len(rule.ExtraActions) > 0
			if !gateHasExtraActions && !gateHasActions {
				lintLog("warning: rule %v will be ignored since it has no actions", ruleName)
			}
		}

		if !gateHasActions && !gateHasExtraActions {
			lintLog("warning: gate has no actions")
		}

		protectionGatesName = append(protectionGatesName, gate.Name)
	}

	return nil
}

// Validations
// - Check that all rules are being used on gates
func lintRulesWithGates(rules map[string]PadRule, padGates []PadProtectionGate) error {
	totalUsesByRule := make(map[string]int, len(rules))

	for ruleName := range rules {
		totalUsesByRule[ruleName] = 0
	}

	for _, gate := range padGates {
		for _, rule := range gate.PatchRules {
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

func lintGroupsMentions(groups map[string]PadGroup, rules map[string]PadRule, gates []PadProtectionGate) error {
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

	for _, gate := range gates {
		actions := gate.Actions
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

	err = lintGates(file.Rules, file.ProtectionGates)
	if err != nil {
		return err
	}

	err = lintRulesWithGates(file.Rules, file.ProtectionGates)
	if err != nil {
		return err
	}

	return lintGroupsMentions(file.Groups, file.Rules, file.ProtectionGates)
}
