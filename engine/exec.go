// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/reviewpad/reviewpad/v4/engine/commands"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/sirupsen/logrus"
)

// EvalCommand generates the program to be executed when a command is received
func EvalCommand(command string, env *Env) (*Program, error) {
	program := BuildProgram(make([]*Statement, 0))

	out := new(bytes.Buffer)
	command = strings.TrimPrefix(command, "/reviewpad")

	args, err := shellwords.Parse(command)
	if err != nil {
		return nil, err
	}

	root := commands.NewCommands(out, args)

	err = root.Execute()
	if err != nil {
		return nil, err
	}

	program.append([]string{out.String()})

	return program, nil
}

// EvalConfigurationFile generates the program to be executed
// Pre-condition Lint(file) == nil
func EvalConfigurationFile(file *ReviewpadFile, env *Env) (*Program, error) {
	log := env.Logger
	interpreter := env.Interpreter

	rules := make(map[string]PadRule)

	log.WithFields(logrus.Fields{
		"reviewpad_file": file,
		"reviewpad_details": map[string]interface{}{
			"project":        fmt.Sprintf(env.TargetEntity.Owner + "/" + env.TargetEntity.Repo),
			"mode":           file.Mode,
			"ignoreErrors":   file.IgnoreErrors,
			"metricsOnMerge": file.MetricsOnMerge,
			"totalImports":   len(file.Imports),
			"totalExtends":   len(file.Extends),
			"totalGroups":    len(file.Groups),
			"totalLabels":    len(file.Labels),
			"totalRules":     len(file.Rules),
			"totalWorkflows": len(file.Workflows),
			"totalPipelines": len(file.Pipelines),
			"totalRecipes":   len(file.Recipes),
		},
	}).Debugln("reviewpad file")

	// process labels
	for labelKeyName, label := range file.Labels {
		labelName := labelKeyName
		// for backwards compatibility, a label has both a key and a name
		if label.Name != "" {
			labelName = label.Name
		} else {
			label.Name = labelName
		}

		if !env.DryRun {
			ghLabel, err := checkLabelExists(env, labelName)
			if err != nil {
				return nil, err
			}

			if ghLabel != nil {
				labelHasUpdates, errCheckUpdates := checkLabelHasUpdates(env, &label, ghLabel)
				if errCheckUpdates != nil {
					return nil, errCheckUpdates
				}

				if labelHasUpdates {
					err = updateLabel(env, &labelName, &label)
					if err != nil {
						return nil, err
					}
				}
			} else {
				err = createLabel(env, &labelName, &label)
				if err != nil {
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
		err := interpreter.ProcessGroup(
			group.Name,
			GroupKind(group.Kind),
			GroupType(group.Type),
			group.Spec,
			group.Param,
			transformAladinoExpression(group.Where),
		)
		if err != nil {
			return nil, err
		}
	}

	// a program is a list of statements to be executed based on the command, workflow rules and actions.
	program := BuildProgram(make([]*Statement, 0))

	// process rules
	for _, rule := range file.Rules {
		err := interpreter.ProcessRule(rule.Name, rule.Spec)
		if err != nil {
			return nil, err
		}
		rules[rule.Name] = rule
	}

	// triggeredExclusiveWorkflow is a control variable to denote if a workflow `always-run: false` has been triggered.
	triggeredExclusiveWorkflow := false

	// process workflows
	for _, workflow := range file.Workflows {
		log.Infof("evaluating workflow '%v'", workflow.Name)
		workflowLog := log.WithField("workflow", workflow.Name)

		if (!workflow.AlwaysRun && !workflow.RunUsed) && triggeredExclusiveWorkflow {
			workflowLog.Infof("skipping workflow because it is not always run and another workflow has been triggered")
			continue
		}

		shouldRun := false
		for _, eventKind := range workflow.On {
			if eventKind == env.TargetEntity.Kind {
				shouldRun = true
				break
			}
		}

		if !shouldRun {
			workflowLog.Infof("skipping workflow because event kind is '%v' and workflow is on '%v'", env.TargetEntity.Kind, workflow.On)
			continue
		}

		for _, run := range workflow.Runs {
			runActions, err := getActionsFromRunBlock(interpreter, &run, rules)
			if err != nil {
				return nil, err
			}

			if len(runActions) > 0 {
				if !workflow.AlwaysRun {
					triggeredExclusiveWorkflow = true
				}

				program.append(runActions)
			}
		}
	}

	// pipelines should only run on pull requests
	if env.TargetEntity.Kind == handler.PullRequest {
		// process pipelines
		for _, pipeline := range file.Pipelines {
			log.Infof("evaluating pipeline '%v':", pipeline.Name)
			pipelineLog := log.WithField("pipeline", pipeline.Name)

			var err error
			activated := pipeline.Trigger == ""
			if !activated {
				activated, err = interpreter.EvalExpr("patch", pipeline.Trigger)
				if err != nil {
					return nil, err
				}
			}

			if activated {
				for num, stage := range pipeline.Stages {
					pipelineLog.Infof("evaluating pipeline stage '%v'", num)
					if stage.Until == "" {
						program.append(stage.Actions)
						break
					}

					isDone, err := interpreter.EvalExpr("patch", stage.Until)
					if err != nil {
						return nil, err
					}

					if !isDone {
						program.append(stage.Actions)
						break
					}
				}
			}
		}
	}

	return program, nil
}

func getActionsFromRunBlock(interpreter Interpreter, run *PadWorkflowRunBlock, rules map[string]PadRule) ([]string, error) {
	if run == nil {
		return nil, nil
	}

	// if the run block was just a simple string
	// there is no rule to evaluate, so just return the actions
	if run.If == nil {
		return run.Actions, nil
	}

	for _, rule := range run.If {
		ruleName := rule.Rule
		ruleDefinition := rules[ruleName]

		activated, err := interpreter.EvalExpr(ruleDefinition.Kind, ruleDefinition.Spec)
		if err != nil {
			return nil, err
		}

		if activated {
			if len(run.Then) > 0 {
				return getActionsFromRunBlocks(interpreter, run.Then, rules, rule.ExtraActions)
			}

			return append(run.Actions, rule.ExtraActions...), nil
		}

		if run.Else != nil {
			return getActionsFromRunBlocks(interpreter, run.Else, rules, nil)
		}
	}

	return nil, nil
}

func getActionsFromRunBlocks(interpreter Interpreter, runs []PadWorkflowRunBlock, rules map[string]PadRule, extraActions []string) ([]string, error) {
	actions := []string{}
	for _, run := range runs {
		runActions, err := getActionsFromRunBlock(interpreter, &run, rules)
		if err != nil {
			return nil, err
		}

		actions = append(actions, append(runActions, extraActions...)...)
	}

	return actions, nil
}
