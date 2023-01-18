// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"
	"strings"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func DryRun() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           dryRunCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func dryRunCode(e aladino.Env, args []aladino.Value) error {
	configFile := e.GetConfigFile()
	if configFile == nil {
		return nil
	}

	program, err := buildProgram(e, configFile)
	if err != nil {
		return err
	}

	report, err := buildReport(e, program)
	if err != nil {
		return err
	}

	return aladino.AddReportComment(e, report)
}

func buildProgram(e aladino.Env, file *engine.ReviewpadFile) (*engine.Program, error) {
	// a program is a list of statements to be executed based on the command, workflow rules and actions.
	program := engine.BuildProgram(make([]*engine.Statement, 0))

	rules := make(map[string]engine.PadRule)
	for _, rule := range file.Rules {
		rules[rule.Name] = rule
	}

	// triggeredExclusiveWorkflow is a control variable to denote if a workflow `always-run: false` has been triggered.
	triggeredExclusiveWorkflow := false

	for _, workflow := range file.Workflows {
		if !workflow.AlwaysRun && triggeredExclusiveWorkflow {
			continue
		}

		ruleActivatedQueue := make([]engine.PadWorkflowRule, 0)
		ruleDefinitionQueue := make(map[string]engine.PadRule)

		shouldRun := false
		for _, eventKind := range workflow.On {
			if eventKind == e.GetTarget().GetTargetEntity().Kind {
				shouldRun = true
				break
			}
		}

		if !shouldRun {
			continue
		}

		for _, rule := range workflow.Rules {
			ruleName := rule.Rule
			ruleDefinition := rules[ruleName]

			activated, err := aladino.EvalExpr(e, ruleDefinition.Kind, ruleDefinition.Spec)
			if err != nil {
				return nil, err
			}

			if activated {
				ruleActivatedQueue = append(ruleActivatedQueue, rule)
				ruleDefinitionQueue[ruleName] = ruleDefinition
			}
		}

		if len(ruleActivatedQueue) > 0 {
			program.Append(workflow.Actions)

			for _, activatedRule := range ruleActivatedQueue {
				program.Append(activatedRule.ExtraActions)
			}

			if !workflow.AlwaysRun {
				triggeredExclusiveWorkflow = true
			}
		}
	}

	for _, pipeline := range file.Pipelines {
		var err error
		activated := pipeline.Trigger == ""
		if !activated {
			activated, err = aladino.EvalExpr(e, "patch", pipeline.Trigger)
			if err != nil {
				return nil, err
			}
		}

		if activated {
			for _, stage := range pipeline.Stages {
				if stage.Until == "" {
					program.Append(stage.Actions)
					break
				}

				isDone, err := aladino.EvalExpr(e, "patch", stage.Until)
				if err != nil {
					return nil, err
				}

				if !isDone {
					program.Append(stage.Actions)
					break
				}
			}
		}
	}

	return program, nil
}

func buildReport(e aladino.Env, program *engine.Program) (string, error) {
	var report strings.Builder
	executedActions := []string{}

	for _, statement := range program.GetProgramStatements() {
		statRaw := statement.GetStatementCode()
		statAST, err := aladino.Parse(statRaw)
		if err != nil {
			return "", err
		}

		_, err = aladino.TypeCheckExec(e, statAST)
		if err != nil {
			return "", err
		}

		executedActions = append(executedActions, statement.GetStatementCode())
	}

	// Annotation
	report.WriteString("<!--@annotation-reviewpad-report-triggered-by-dry-run-command-->")
	// Header
	report.WriteString("**Reviewpad Report** (Reviewpad ran in dry-run mode because the `$dryRun() action was triggered`)\n\n")

	report.WriteString(":scroll: **Executed actions**\n")
	report.WriteString("```yaml\n")

	// Report
	for _, action := range executedActions {
		report.WriteString(fmt.Sprintf("%v\n", action))
	}

	report.WriteString("```\n")

	return report.String(), nil
}
