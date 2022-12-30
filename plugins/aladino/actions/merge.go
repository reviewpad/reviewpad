// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Merge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           mergeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func mergeCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)

	if t.PullRequest.GetState() != "open" || t.PullRequest.GetDraft() {
		e.GetLogger().Infof("merge: skipping action because pull request is not open or is a draft")
		return nil
	}

	mergeMethod, err := parseMergeMethod(args)
	if err != nil {
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
