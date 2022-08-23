// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Milestone() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildStringType()),
		Code:           milestoneCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func milestoneCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	milestoneTitle := pullRequest.GetMilestone().GetTitle()
	return aladino.BuildStringValue(milestoneTitle), nil
}
