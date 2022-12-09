// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package reviewpad

import (
	"bytes"
	"context"
	"fmt"
	"log"

	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
)

func Load(ctx context.Context, githubClient *gh.GithubClient, buf *bytes.Buffer) (*engine.ReviewpadFile, error) {
	file, err := engine.Load(ctx, githubClient, buf.Bytes())
	if err != nil {
		return nil, err
	}

	log.Println(fmtio.Sprintf("load", "input file:\n%+v\n", file))

	err = engine.Lint(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func Run(
	ctx context.Context,
	githubClient *gh.GithubClient,
	collector collector.Collector,
	targetEntity *handler.TargetEntity,
	eventData *handler.EventData,
	eventPayload interface{},
	reviewpadFile *engine.ReviewpadFile,
	dryRun bool,
	safeMode bool,
) (engine.ExitStatus, *engine.Program, error) {
	if safeMode && !dryRun {
		return engine.ExitStatusFailure, nil, fmt.Errorf("when reviewpad is running in safe mode, it must also run in dry-run")
	}

	config, err := plugins_aladino.DefaultPluginConfig()
	if err != nil {
		return engine.ExitStatusFailure, nil, err
	}

	defer config.CleanupPluginConfig()

	aladinoInterpreter, err := aladino.NewInterpreter(ctx, dryRun, githubClient, collector, targetEntity, eventPayload, plugins_aladino.PluginBuiltInsWithConfig(config))
	if err != nil {
		return engine.ExitStatusFailure, nil, err
	}

	evalEnv, err := engine.NewEvalEnv(ctx, dryRun, githubClient, collector, targetEntity, aladinoInterpreter, eventData)
	if err != nil {
		return engine.ExitStatusFailure, nil, err
	}

	program, err := engine.Eval(reviewpadFile, evalEnv)
	if err != nil {
		return engine.ExitStatusFailure, nil, err
	}

	exitStatus, err := aladinoInterpreter.ExecProgram(program)
	if err != nil {
		engine.CollectError(evalEnv, err)
		return engine.ExitStatusFailure, nil, err
	}

	if safeMode || !dryRun {
		err = aladinoInterpreter.Report(reviewpadFile.Mode, safeMode)
		if err != nil {
			engine.CollectError(evalEnv, err)
			return engine.ExitStatusFailure, nil, err
		}
	}

	if utils.IsPullRequestReadyForReportMetrics(eventData) {
		err = aladinoInterpreter.ReportMetrics(reviewpadFile.Mode)
		if err != nil {
			engine.CollectError(evalEnv, err)
			return engine.ExitStatusFailure, nil, err
		}
	}

	collectedData := map[string]interface{}{}

	err = evalEnv.Collector.Collect("Completed Analysis", collectedData)

	if err != nil {
		log.Printf("error on collector due to %v", err.Error())
	}

	return exitStatus, program, nil
}
