// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package reviewpad

import (
	"bytes"
	"context"
	"log"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/collector"
	"github.com/reviewpad/reviewpad/engine"
	"github.com/reviewpad/reviewpad/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/plugins/aladino"
	"github.com/reviewpad/reviewpad/utils/fmtio"
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
	ghPullRequest *github.PullRequest,
	reviewpadFile *engine.ReviewpadFile,
	dryRun bool,
) error {
	interpreters := make(map[string]engine.Interpreter)

	aladinoInterpreter, err := aladino.NewInterpreter(ctx, client, clientGQL, collector, ghPullRequest, plugins_aladino.PluginBuiltIns())
	if err != nil {
		return err
	}

	interpreters["aladino"] = aladinoInterpreter
	evalEnv, err := engine.NewEvalEnv(ctx, client, clientGQL, collector, ghPullRequest, interpreters)
	if err != nil {
		return err
	}

	_, err = engine.Exec(reviewpadFile, evalEnv, &engine.Flags{Dryrun: dryRun})
	if err != nil {
		return err
	}

	return nil
}
