// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func Merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: mergeCode,
	}
}

func mergeCode(e aladino.Env, args []aladino.Value) error {
	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestRepoName(e.GetPullRequest())

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
		return err
	}

	_, _, err = e.GetClient().PullRequests.Merge(e.GetCtx(), owner, repo, prNum, "Merged by Reviewpad", &github.PullRequestOptions{
		MergeMethod: mergeMethod,
	})
	return err
}

func parseMergeMethod(args []aladino.Value) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf("merge: received two arguments")
	}

	if len(args) == 0 {
		return "merge", nil
	}

	arg := args[0]
	if arg.HasKindOf(aladino.STRING_VALUE) {
		mergeMethod := arg.(*aladino.StringValue).Val
		switch mergeMethod {
		case "merge", "rebase", "squash":
			return mergeMethod, nil
		default:
			return "", fmt.Errorf("merge: unexpected argument %v", mergeMethod)
		}
	} else {
		return "", fmt.Errorf("merge: expects string argument")
	}
}
