// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"log"

	"github.com/reviewpad/reviewpad/utils/fmtio"
)

func execError(format string, a ...interface{}) error {
	return fmtio.Errorf("reviewpad", format, a...)
}

func execLogf(format string, a ...interface{}) {
	log.Println(fmtio.Sprintf("reviewpad", format, a...))
}

func execLog(val string) {
	log.Println(fmtio.Sprint("reviewpad", val))
}

// Exec: main function
// Pre-condition Lint(file) == nil
func Exec(file *ReviewpadFile, env *Env, flags *Flags) ([]string, error) {
	execLogf("file to execute:\n%+v", file)

	rules := make(map[string]PadRule)

	reportDetails := make([]ReportGateDetails, 0)

	execLogf("detected %v groups", len(file.Groups))
	execLogf("detected %v labels", len(file.Labels))
	execLogf("detected %v rules", len(file.Rules))
	execLogf("detected %v gates", len(file.ProtectionGates))

	// process labels
	for labelName, label := range file.Labels {
		labelInRepo, _ := getLabel(env, &labelName)

		if labelInRepo.GetName() != "" {
			// label already exists: nothing to do
			continue
		}

		err := createLabel(env, &labelName, &label)
		if err != nil {
			return nil, err
		}
	}

	// lang := file.Language
	interpreter, ok := env.Interpreters["aladino"]
	if !ok {
		return nil, execError("no interpreter for aladino")
	}

	// process groups
	for groupName, group := range file.Groups {
		err := interpreter.ProcessGroup(groupName, GroupKind(group.Kind), GroupType(group.Type), group.Spec, group.Param, group.Where)
		if err != nil {
			return nil, err
		}
	}

	for ruleName, rule := range file.Rules {
		rules[ruleName] = rule
	}

	// program defines the set of all actions to run.
	// This set is calculated based on the gates rules and actions.
	program := make([]string, 0)
	// triggeredGate defines the gate `alwaysRun: false` that has been triggered.
	triggeredGate := ""

	for _, gate := range file.ProtectionGates {
		execLogf("evaluating gate %v:", gate.Name)

		isGateEligible := gate.AlwaysRun || triggeredGate == ""

		if !isGateEligible {
			execLog("\tnot eligible")
			continue
		}

		reportGateDetails := ReportGateDetails{
			name:            gate.Name,
			description:     gate.Description,
			actRules:        []string{},
			actActions:      []string{},
			actExtraActions: []string{},
		}

		ruleActivatedQueue := make([]PatchRule, 0)
		ruleDefinitionQueue := make(map[string]PadRule)

		for _, rule := range gate.PatchRules {
			ruleName := rule.Rule
			ruleDefinition := rules[ruleName]

			activated, err := interpreter.EvalExpr(ruleDefinition.Kind, ruleDefinition.Spec)
			if err != nil {
				return nil, err
			}

			if activated {
				ruleActivatedQueue = append(ruleActivatedQueue, rule)
				ruleDefinitionQueue[ruleName] = ruleDefinition

				execLogf("\trule %v activated", ruleName)
			}
		}

		if len(ruleActivatedQueue) > 0 && isGateEligible {
			program = append(program, gate.Actions...)

			reportGateDetails.actActions = append(reportGateDetails.actActions, gate.Actions...)

			for _, activatedRule := range ruleActivatedQueue {
				program = append(program, activatedRule.ExtraActions...)

				reportGateDetails.actRules = append(reportGateDetails.actRules, activatedRule.Rule)
				reportGateDetails.actExtraActions = append(reportGateDetails.actExtraActions, activatedRule.ExtraActions...)
			}

			if !gate.AlwaysRun {
				triggeredGate = gate.Name
			}

			reportDetails = append(reportDetails, reportGateDetails)
		} else {
			execLog("\tno rules activated")
		}
	}

	if !flags.Dryrun {
		if file.Mode != "silent" {
			// Although reportProgram is async we don't wait for it to finish.
			// Also, we don't care if the report was done successfully or not.
			reportProgram(env, &reportDetails)
		}

		err := interpreter.ExecActions(&program)
		if err != nil {
			return nil, err
		}
	}

	execLogf("executed program:\n%+q", program)

	return program, nil
}
