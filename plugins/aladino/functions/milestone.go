// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Milestone() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildStringType()),
		Code:           milestoneCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func milestoneCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget).PullRequest
	milestoneTitle := pullRequest.GetMilestone().GetTitle()
	return lang.BuildStringValue(milestoneTitle), nil
}
