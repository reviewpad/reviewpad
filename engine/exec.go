// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"log"
	"regexp"

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

func collectError(env *Env, err error) {
	env.Collector.Collect("Error", &map[string]interface{}{
		"pullRequestUrl": env.PullRequest.URL,
		"details":        err.Error(),
	})
}

// Exec: main function
// Pre-condition Lint(file) == nil
func Exec(file *ReviewpadFile, env *Env, flags *Flags) ([]string, error) {
	execLogf("file to execute:\n%+v", file)

	reg := regexp.MustCompile(`github\.com\/repos\/(.*)\/pulls\/\d+$$`)
	matches := reg.FindStringSubmatch(*env.PullRequest.URL)

	env.Collector.Collect("Trigger Analysis", &map[string]interface{}{
		"pullRequestUrl": env.PullRequest.URL,
		"project":        matches[1],
		"version":        file.Version,
		"edition":        file.Edition,
		"mode":           file.Mode,
		"totalGroups":    len(file.Groups),
		"totalLabels":    len(file.Labels),
		"totalRules":     len(file.Rules),
		"totalWorkflows": len(file.Workflows),
	})

	rules := make(map[string]PadRule)

	reportDetails := make([]ReportWorkflowDetails, 0)

	execLogf("detected %v groups", len(file.Groups))
	execLogf("detected %v labels", len(file.Labels))
	execLogf("detected %v rules", len(file.Rules))
	execLogf("detected %v workflows", len(file.Workflows))

	// process labels
	for labelName, label := range file.Labels {
		labelInRepo, _ := getLabel(env, &labelName)

		if labelInRepo.GetName() != "" {
			// label already exists: nothing to do
			continue
		}

		err := createLabel(env, &labelName, &label)
		if err != nil {
			collectError(env, err)
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
			collectError(env, err)
			return nil, err
		}
	}

	for ruleName, rule := range file.Rules {
		rules[ruleName] = rule
	}

	// program defines the set of all actions to run.
	// This set is calculated based on the workflow rules and actions.
	program := make([]string, 0)
	// triggeredWorkflow defines the workflow `always-run: false` that has been triggered.
	triggeredWorkflow := ""

	for _, workflow := range file.Workflows {
		execLogf("evaluating workflow %v:", workflow.Name)

		isWorkflowEligible := workflow.AlwaysRun || triggeredWorkflow == ""

		if !isWorkflowEligible {
			execLog("\tnot eligible")
			continue
		}

		reportWorkflowDetails := ReportWorkflowDetails{
			name:            workflow.Name,
			description:     workflow.Description,
			actRules:        []string{},
			actActions:      []string{},
			actExtraActions: []string{},
		}

		ruleActivatedQueue := make([]PatchRule, 0)
		ruleDefinitionQueue := make(map[string]PadRule)

		for _, rule := range workflow.PatchRules {
			ruleName := rule.Rule
			ruleDefinition := rules[ruleName]

			activated, err := interpreter.EvalExpr(ruleDefinition.Kind, ruleDefinition.Spec)
			if err != nil {
				collectError(env, err)
				return nil, err
			}

			if activated {
				ruleActivatedQueue = append(ruleActivatedQueue, rule)
				ruleDefinitionQueue[ruleName] = ruleDefinition

				execLogf("\trule %v activated", ruleName)
			}
		}

		if len(ruleActivatedQueue) > 0 && isWorkflowEligible {
			program = append(program, workflow.Actions...)

			reportWorkflowDetails.actActions = append(reportWorkflowDetails.actActions, workflow.Actions...)

			for _, activatedRule := range ruleActivatedQueue {
				program = append(program, activatedRule.ExtraActions...)

				reportWorkflowDetails.actRules = append(reportWorkflowDetails.actRules, activatedRule.Rule)
				reportWorkflowDetails.actExtraActions = append(reportWorkflowDetails.actExtraActions, activatedRule.ExtraActions...)
			}

			if !workflow.AlwaysRun {
				triggeredWorkflow = workflow.Name
			}

			reportDetails = append(reportDetails, reportWorkflowDetails)
		} else {
			execLog("\tno rules activated")
		}
	}

	if !flags.Dryrun {
		if file.Mode == VERBOSE_MODE {
			// Ignore errors in the report of the verbose mode
			reportProgram(env, &reportDetails)
		}

		err := interpreter.ExecActions(&program)
		if err != nil {
			collectError(env, err)
			return nil, err
		}
	}

	execLogf("executed program:\n%+q", program)

	env.Collector.Collect("Completed Analysis", &map[string]interface{}{
		"pullRequestUrl": env.PullRequest.URL,
	})

	return program, nil
}
