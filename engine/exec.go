// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"log"
	"regexp"

	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
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

func CollectError(env *Env, err error) {
	env.Collector.Collect("Error", map[string]interface{}{
		"pullRequestUrl": env.PullRequest.URL,
		"details":        err.Error(),
	})
}

// Eval: main function that generates the program to be executed
// Pre-condition Lint(file) == nil
func Eval(file *ReviewpadFile, env *Env, reportSettings ReportSettings) (*Program, error) {
	execLogf("file to evaluate:\n%+v", file)

	interpreter := env.Interpreter

	reg := regexp.MustCompile(`github\.com\/repos\/(.*)\/pulls\/\d+$$`)
	matches := reg.FindStringSubmatch(*env.PullRequest.URL)

	env.Collector.Collect("Trigger Analysis", map[string]interface{}{
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

	execLogf("detected %v groups", len(file.Groups))
	execLogf("detected %v labels", len(file.Labels))
	execLogf("detected %v rules", len(file.Rules))
	execLogf("detected %v workflows", len(file.Workflows))

	// process labels
	for labelKeyName, label := range file.Labels {
		labelName := labelKeyName
		// for backwards compatibility, a label has both a key and a name
		if label.Name != "" {
			labelName = label.Name
		}

		if !env.DryRun {
			labelExists, err := checkLabelExists(env, labelName)
			if err != nil {
				return nil, err
			}

			if !labelExists {
				err = createLabel(env, &labelName, &label)
				if err != nil {
					CollectError(env, err)
					return nil, err
				}
			}
		}

		err := interpreter.ProcessLabel(labelKeyName, labelName)
		if err != nil {
			return nil, err
		}
	}

	// process groups
	for _, group := range file.Groups {
		err := interpreter.ProcessGroup(group.Name, GroupKind(group.Kind), GroupType(group.Type), group.Spec, group.Param, group.Where)
		if err != nil {
			CollectError(env, err)
			return nil, err
		}
	}

	// process rules
	for _, rule := range file.Rules {
		err := interpreter.ProcessRule(rule.Name, rule.Spec)
		if err != nil {
			CollectError(env, err)
			return nil, err
		}
		rules[rule.Name] = rule
	}

	// a program is a list of statements to be executed based on the workflow rules and actions.
	program := &Program{
		ReportSettings: reportSettings,
		Statements:     make([]*Statement, 0),
	}

	// triggeredExclusiveWorkflow is a control variable to denote if a workflow `always-run: false` has been triggered.
	triggeredExclusiveWorkflow := false

	for _, workflow := range file.Workflows {
		execLogf("evaluating workflow %v:", workflow.Name)

		if !workflow.AlwaysRun && triggeredExclusiveWorkflow {
			execLog("\tskipping workflow")
			continue
		}

		ruleActivatedQueue := make([]PadWorkflowRule, 0)
		ruleDefinitionQueue := make(map[string]PadRule)

		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			ruleDefinition := rules[ruleName]

			activated, err := interpreter.EvalExpr(ruleDefinition.Kind, ruleDefinition.Spec)
			if err != nil {
				CollectError(env, err)
				return nil, err
			}

			if activated {
				ruleActivatedQueue = append(ruleActivatedQueue, rule)
				ruleDefinitionQueue[ruleName] = ruleDefinition

				execLogf("\trule %v activated", ruleName)
			}
		}

		if len(ruleActivatedQueue) > 0 {
			program.append(workflow.Actions, workflow, ruleActivatedQueue)

			for _, activatedRule := range ruleActivatedQueue {
				program.append(activatedRule.ExtraActions, workflow, []PadWorkflowRule{activatedRule})
			}

			if !workflow.AlwaysRun {
				triggeredExclusiveWorkflow = true
			}
		} else {
			execLog("\tno rules activated")
		}
	}

	return program, nil
}
