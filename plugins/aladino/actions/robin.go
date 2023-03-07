// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/robin"
)

func Robin() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           robinCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func robinCode(e aladino.Env, args []aladino.Value) error {
	target := e.GetTarget().(*target.PullRequestTarget)
	prompt := args[0].(*aladino.StringValue).Val

	reply, err := robin.Prompt(e, prompt)
	if err != nil {
		return nil
	}

	return target.Comment(reply)
}
