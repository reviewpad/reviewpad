// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func PullRequestCountBy() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildIntType()),
		Code:           pullRequestCountByCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func pullRequestCountByCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	loginArg := args[0].(*lang.StringValue).Val
	stateArg := args[1].(*lang.StringValue).Val

	opts := &github.IssueListByRepoOptions{
		State: stateArg,
	}
	if loginArg != "" {
		opts.Creator = loginArg
	}

	return issueCountBy(e, opts, func(issue *github.Issue) bool {
		return issue.IsPullRequest()
	})
}
