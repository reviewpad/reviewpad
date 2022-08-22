// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func LastEventAt() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code: lastEventAtCode,
	}
}

func lastEventAtCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetPullRequest()
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)
	prNum := gh.GetPullRequestNumber(pullRequest)

	timeline, err := e.GetGithubClient().GetIssueTimeline(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	lastEvent := timeline[len(timeline)-1]
	var lastEventTime int
	if *lastEvent.Event == "reviewed" {
		lastEventTime = int(lastEvent.GetSubmittedAt().Unix())
	} else {
		lastEventTime = int(lastEvent.GetCreatedAt().Unix())
	}

	return aladino.BuildIntValue(lastEventTime), nil
}
