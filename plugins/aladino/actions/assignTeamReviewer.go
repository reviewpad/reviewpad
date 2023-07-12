// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func AssignTeamReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildArrayOfType(lang.BuildStringType())}, nil),
		Code:           assignTeamReviewerCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func assignTeamReviewerCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)

	teamReviewers := args[0].(*lang.ArrayValue).Vals

	if len(teamReviewers) < 1 {
		return fmt.Errorf("assignTeamReviewer: requires at least 1 team to request for review")
	}

	teamReviewersSlugs := make([]string, len(teamReviewers))

	for i, team := range teamReviewers {
		teamReviewersSlugs[i] = team.(*lang.StringValue).Val
	}

	return t.RequestTeamReviewers(teamReviewersSlugs)
}
