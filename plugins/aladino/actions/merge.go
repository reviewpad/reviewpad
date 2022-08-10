// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/google/go-github/v45/github"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: mergeCode,
	}
}

func mergeCode(e aladino.Env, args []aladino.Value) error {
	pullRequest := e.GetPullRequest()
	prNum := gh.GetPullRequestNumber(pullRequest)
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
		return err
	}

	_, _, err = e.GetGithubClient().Merge(e.GetCtx(), owner, repo, prNum, "Merged by Reviewpad", &github.PullRequestOptions{
		MergeMethod: mergeMethod,
	})
	return err
}

func parseMergeMethod(args []aladino.Value) (string, error) {
	if len(args) == 0 {
		return "merge", nil
	}

	mergeMethod := args[0].(*aladino.StringValue).Val
	switch mergeMethod {
	case "merge", "rebase", "squash":
		return mergeMethod, nil
	default:
		return "", fmt.Errorf("merge: unsupported merge method %v", mergeMethod)
	}
}
