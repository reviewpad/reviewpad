// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/google/go-github/v52/github"
	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           mergeCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func mergeCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	targetEntity := e.GetTarget().GetTargetEntity()
	log := e.GetLogger().WithField("builtin", "merge")

	if t.PullRequest.Status != pbc.PullRequestStatus_OPEN || t.PullRequest.IsDraft {
		log.Infof("skipping action because pull request is not open or is a draft")
		return nil
	}

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
		return err
	}

	if e.GetCheckRunID() != nil {
		e.SetCheckRunUpdated(true)
		_, _, err := e.GetGithubClient().GetClientREST().Checks.UpdateCheckRun(e.GetCtx(), targetEntity.Owner, targetEntity.Repo, *e.GetCheckRunID(), github.UpdateCheckRunOptions{
			Status:     github.String("completed"),
			Conclusion: github.String("success"),
		})
		return err
	}

	return t.Merge(mergeMethod)
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
