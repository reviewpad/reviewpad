// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func DisableMerge() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           disableMergeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func disableMergeCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)

	_, err := t.CreateCommitStatus(t.PullRequest.GetHead().GetSHA(), &github.CreateCommitStatusOptions{
		Context:     "reviewpad merge gate",
		State:       "failure",
		Description: args[0].(*aladino.StringValue).Val,
	})

	return err
}
