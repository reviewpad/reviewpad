// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package reviewpad

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/reviewpad/reviewpad/v3/utils/fmtio"
	"github.com/shurcooL/githubv4"
)

func Load(buf *bytes.Buffer) (*engine.ReviewpadFile, error) {
	file, err := engine.Load(buf.Bytes())
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
	client *github.Client,
	clientGQL *githubv4.Client,
	collector collector.Collector,
	pullRequest *github.PullRequest,
	eventPayload interface{},
	reviewpadFile *engine.ReviewpadFile,
	dryRun bool,
	safeMode bool,
) (*engine.Program, error) {
	if safeMode && !dryRun {
		return nil, fmt.Errorf("when reviewpad is running in safe mode, it must also run in dry-run")
	}

	aladinoInterpreter, err := aladino.NewInterpreter(ctx, dryRun, client, clientGQL, collector, pullRequest, eventPayload, plugins_aladino.PluginBuiltIns())
	if err != nil {
		return nil, err
	}

	evalEnv, err := engine.NewEvalEnv(ctx, dryRun, client, clientGQL, collector, pullRequest, eventPayload, aladinoInterpreter)
	if err != nil {
		return nil, err
	}

	reportSettings := engine.ReportSettings{
		UseSafeModeHeader: safeMode,
		IsReportEnabled:   reviewpadFile.Mode == "VERBOSE",
	}

	program, err := engine.Eval(reviewpadFile, evalEnv, reportSettings)
	if err != nil {
		return nil, err
	}

	err = aladinoInterpreter.ExecProgram(program)
	if err != nil {
		engine.CollectError(evalEnv, err)
		return nil, err
	}

	if safeMode || !dryRun {
		err = aladinoInterpreter.Report(reviewpadFile.Mode, safeMode)
		if err != nil {
			engine.CollectError(evalEnv, err)
			return nil, err
		}
	}

	err = evalEnv.Collector.Collect("Completed Analysis", map[string]interface{}{
		"pullRequestUrl": evalEnv.PullRequest.URL,
	})

	if err != nil {
		log.Printf("error on collector due to %v", err.Error())
	}

	return program, nil
}
