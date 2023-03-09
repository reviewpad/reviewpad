// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Comment() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           commentCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func commentCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	pullRequest := t.PullRequest
	repo := pullRequest.GetBase().GetRepo()
	commentBody := args[0].(*aladino.StringValue).Val

	hostClient := e.GetCodeHostClient()
	return hostClient.PostGeneralComment(e.GetCtx(), repo.GetFullName(), repo.GetNodeID(), pullRequest.GetNodeID(), int32(pullRequest.GetNumber()), commentBody)
}
