// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package reviewpad

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	github_event_processor "github.com/reviewpad/go-lib/event/event_processor/github"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/collector"
	"github.com/reviewpad/reviewpad/v4/engine"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/sirupsen/logrus"
)

func logErrorAndCollect(logger *logrus.Entry, collector collector.Collector, message string, err error) {
	logger.WithError(err).Errorln(fmt.Sprintf("%s `%s`", message, err.Error()))

	if ghError, isGitHubError := err.(*github.ErrorResponse); isGitHubError {
		err = collector.CollectError(fmt.Errorf("%s: %s", message, ghError.Message))
	} else {
		err = collector.CollectError(fmt.Errorf("%s: %s", message, err.Error()))
	}
	if err != nil {
		logger.WithError(err).Errorf("failed to collect error `%s`", err.Error())
	}
}

func Load(ctx context.Context, log *logrus.Entry, githubClient *gh.GithubClient, buf *bytes.Buffer) (*engine.ReviewpadFile, error) {
	file, err := engine.Load(ctx, log, githubClient, buf.Bytes())
	if err != nil {
		return nil, err
	}

	log.WithField("reviewpad_file", file).Debug("loaded reviewpad file")

	builtIns := plugins_aladino.PluginBuiltIns()
	reserved := []string{}
	for name := range builtIns.Actions {
		reserved = append(reserved, name)
	}
	for name := range builtIns.Functions {
		reserved = append(reserved, name)
	}

	err = engine.Lint(file, reserved, log)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createCommentWithReviewpadCommandEvaluationError(env *engine.Env, err error) error {
	command := env.EventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody()

	commentBody := new(strings.Builder)
	commentBody.WriteString(fmt.Sprintf("> %s\n\n", command))
	commentBody.WriteString(fmt.Sprintf("\n*Errors:*\n- %s\n", err.Error()))

	_, _, err = env.GithubClient.CreateComment(env.Ctx, env.TargetEntity.Owner, env.TargetEntity.Repo, env.TargetEntity.Number, &github.IssueComment{
		Body: github.String(commentBody.String()),
	})

	return err
}

func createCommentWithReviewpadCommandExecutionError(env *engine.Env, err error) error {
	command := env.EventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody()
	sender := env.EventDetails.Payload.(*github.IssueCommentEvent).GetSender().GetLogin()

	commentBody := new(strings.Builder)
	commentBody.WriteString(fmt.Sprintf("> %s\n\n", command))
	commentBody.WriteString(fmt.Sprintf("@%s an error occurred running your command\n", sender))
	commentBody.WriteString("\n*Errors:*\n")

	gitHubError := &github.ErrorResponse{}
	if errors.As(err, &gitHubError) {
		if len(gitHubError.Errors) > 0 {
			for _, e := range gitHubError.Errors {
				commentBody.WriteString(fmt.Sprintf("- %s\n", e.Message))
			}
		} else {
			commentBody.WriteString(fmt.Sprintf("- %s\n", gitHubError.Message))
		}
	} else {
		commentBody.WriteString(fmt.Sprintf("- %s\n", err.Error()))
	}

	_, _, err = env.GithubClient.CreateComment(env.Ctx, env.TargetEntity.Owner, env.TargetEntity.Repo, env.TargetEntity.Number, &github.IssueComment{
		Body: github.String(commentBody.String()),
	})

	return err
}

func createCommentWithReviewpadCommandExecutionSuccess(env *engine.Env, executionDetails string) error {
	command := env.EventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody()

	commentBody := new(strings.Builder)
	commentBody.WriteString(fmt.Sprintf("> %s\n\n", command))
	commentBody.WriteString(fmt.Sprintf("\n*Execution details:*\n\n%s", executionDetails))

	_, _, err := env.GithubClient.CreateComment(env.Ctx, env.TargetEntity.Owner, env.TargetEntity.Repo, env.TargetEntity.Number, &github.IssueComment{
		Body: github.String(commentBody.String()),
	})

	return err
}

func runReviewpadCommandDryRun(
	logger *logrus.Entry,
	collector collector.Collector,
	gitHubClient *gh.GithubClient,
	targetEntity *entities.TargetEntity,
	reviewpadFile *engine.ReviewpadFile,
	env *engine.Env,
) (engine.ExitStatus, *engine.Program, string, error) {
	err := collector.Collect("Trigger Command", map[string]interface{}{
		"project": fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo),
		"command": "dry-run",
	})
	if err != nil {
		logger.WithError(err).Errorf("error collecting trigger command")
	}

	if reviewpadFile == nil {
		err = errors.New("trying to run reviewpad dry-run command on a repository without a reviewpad file")
		logger.WithError(err).Errorln("error running reviewpad dry-run command")

		err = createCommentWithReviewpadCommandExecutionError(env, err)
		if err != nil {
			logger.WithError(err).Errorf("error creating comment to report error on command execution")
		}

		return engine.ExitStatusFailure, nil, "", err
	}

	program, err := engine.EvalConfigurationFile(reviewpadFile, env)
	if err != nil {
		logErrorAndCollect(logger, collector, "error evaluating configuration file", err)
		return engine.ExitStatusFailure, nil, "", err
	}

	report := aladino.BuildDryRunExecutionReport(program)

	err = createCommentWithReviewpadCommandExecutionSuccess(env, report)
	if err != nil {
		return engine.ExitStatusFailure, program, "", fmt.Errorf("error on creating report comment %v", err.(*github.ErrorResponse).Message)
	}

	return engine.ExitStatusSuccess, program, "", nil
}

func runReviewpadCommand(
	ctx context.Context,
	logger *logrus.Entry,
	collector collector.Collector,
	gitHubClient *gh.GithubClient,
	targetEntity *entities.TargetEntity,
	env *engine.Env,
	aladinoInterpreter engine.Interpreter,
	command string,
) (engine.ExitStatus, *engine.Program, string, error) {
	err := collector.Collect("Trigger Command", map[string]interface{}{
		"project": fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo),
		"command": command,
	})
	if err != nil {
		logger.WithError(err).Errorf("error collecting trigger command")
	}

	program, err := engine.EvalCommand(command, env)
	if err != nil {
		logErrorAndCollect(logger, collector, "error evaluating command", err)

		err = createCommentWithReviewpadCommandEvaluationError(env, err)
		if err != nil {
			logger.WithError(err).Errorf("error creating comment to report error on command execution")
			return engine.ExitStatusFailure, nil, "", err
		}

		return engine.ExitStatusFailure, nil, "", err
	}

	exitStatus, err := aladinoInterpreter.ExecProgram(program)
	if err != nil {
		logErrorAndCollect(logger, collector, "error executing command", err)

		err = createCommentWithReviewpadCommandExecutionError(env, err)
		if err != nil {
			logger.WithError(err).Errorf("error creating comment to report error on command execution")
		}

		return engine.ExitStatusFailure, nil, "", err
	}

	err = collector.Collect("Completed Command Execution", map[string]interface{}{})
	if err != nil {
		logger.WithError(err).Errorf("error collecting completed command execution")
	}

	return exitStatus, program, "", nil
}

func runReviewpadFile(env *engine.Env, reviewpadFile *engine.ReviewpadFile, safeMode bool) (engine.ExitStatus, *engine.Program, string, error) {
	err := env.Collector.Collect("Trigger Analysis", map[string]interface{}{
		"project":        fmt.Sprintf("%s/%s", env.TargetEntity.Owner, env.TargetEntity.Repo),
		"mode":           reviewpadFile.Mode,
		"author":         github_event_processor.GetEventSender(env.EventDetails.Payload),
		"ignoreErrors":   reviewpadFile.IgnoreErrors,
		"metricsOnMerge": reviewpadFile.MetricsOnMerge,
		"totalImports":   len(reviewpadFile.Imports),
		"totalExtends":   len(reviewpadFile.Extends),
		"totalGroups":    len(reviewpadFile.Groups),
		"totalLabels":    len(reviewpadFile.Labels),
		"totalRules":     len(reviewpadFile.Rules),
		"totalWorkflows": len(reviewpadFile.Workflows),
		"totalPipelines": len(reviewpadFile.Pipelines),
		"totalRecipes":   len(reviewpadFile.Recipes),
	})
	if err != nil {
		env.Logger.WithError(err).Errorf("error collecting trigger analysis")
	}

	exitStatus, program, err := engine.ExecConfigurationFile(env, reviewpadFile)
	if err != nil {
		logErrorAndCollect(env.Logger, env.Collector, "error executing configuration file", err)
		return engine.ExitStatusFailure, nil, env.Interpreter.GetCheckRunConclusion(), err
	}

	if safeMode || !env.DryRun {
		err = env.Interpreter.Report(reviewpadFile.Mode, safeMode)
		if err != nil {
			logErrorAndCollect(env.Logger, env.Collector, "error reporting results", err)
			return engine.ExitStatusFailure, nil, env.Interpreter.GetCheckRunConclusion(), err
		}
	}

	if utils.IsPullRequestReadyForReportMetrics(env.EventDetails) {
		if reviewpadFile.MetricsOnMerge != nil && *reviewpadFile.MetricsOnMerge {
			err = env.Interpreter.ReportMetrics()
			if err != nil {
				logErrorAndCollect(env.Logger, env.Collector, "error reporting metrics", err)
				return engine.ExitStatusFailure, nil, env.Interpreter.GetCheckRunConclusion(), err
			}
		}
	}

	err = env.Collector.Collect("Completed Analysis", map[string]interface{}{})
	if err != nil {
		env.Logger.WithError(err).Errorf("error collecting completed analysis")
	}

	return exitStatus, program, env.Interpreter.GetCheckRunConclusion(), nil
}

func Run(
	ctx context.Context,
	log *logrus.Entry,
	gitHubClient *gh.GithubClient,
	codeHostClient *codehost.CodeHostClient,
	collector collector.Collector,
	targetEntity *entities.TargetEntity,
	eventDetails *entities.EventDetails,
	reviewpadFile *engine.ReviewpadFile,
	checkRunId *int64,
	dryRun bool,
	safeMode bool,
	checksWithIssue []string,
) (engine.ExitStatus, *engine.Program, string, error) {
	config, err := plugins_aladino.DefaultPluginConfig()
	if err != nil {
		return engine.ExitStatusFailure, nil, "", err
	}

	defer config.CleanupPluginConfig()

	aladinoInterpreter, err := aladino.NewInterpreter(ctx, log, dryRun, gitHubClient, codeHostClient, collector, targetEntity, eventDetails.Payload, plugins_aladino.PluginBuiltInsWithConfig(config), checkRunId, checksWithIssue)
	if err != nil {
		return engine.ExitStatusFailure, nil, "", err
	}

	env, err := engine.NewEvalEnv(ctx, log, dryRun, gitHubClient, collector, targetEntity, aladinoInterpreter, eventDetails)
	if err != nil {
		return engine.ExitStatusFailure, nil, "", err
	}

	if utils.IsReviewpadCommand(env.EventDetails) {
		// reviewpad-an: fixme: this is specific to github and only works for issues
		command := eventDetails.Payload.(*github.IssueCommentEvent).GetComment().GetBody()
		if utils.IsReviewpadCommandDryRun(command) {
			return runReviewpadCommandDryRun(log, collector, gitHubClient, targetEntity, reviewpadFile, env)
		}

		if utils.IsReviewpadCommandRun(command) {
			return runReviewpadFile(env, reviewpadFile, safeMode)
		}

		return runReviewpadCommand(ctx, log, collector, gitHubClient, targetEntity, env, aladinoInterpreter, command)
	} else {
		return runReviewpadFile(env, reviewpadFile, safeMode)
	}
}
