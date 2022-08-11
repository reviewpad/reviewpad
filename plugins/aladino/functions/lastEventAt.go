// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"time"

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
	owner := gh.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := gh.GetPullRequestBaseRepoName(e.GetPullRequest())
	prNum := gh.GetPullRequestNumber(e.GetPullRequest())

	timeline, err := e.GetGithubClient().GetIssueTimeline(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return nil, err
	}

	lastEventAt := timeline[len(timeline)-1].GetCreatedAt()

	lastEventAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", lastEventAt.String())
	if err != nil {
		return nil, err
	}

	return aladino.BuildIntValue(int(lastEventAtTime.Unix())), nil
}
