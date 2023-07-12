// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v52/github"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func IssueCountBy() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, lang.BuildIntType()),
		Code:           issueCountByCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func issueCountByCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	loginArg := args[0].(*lang.StringValue).Val
	stateArg := args[1].(*lang.StringValue).Val

	opts := &github.IssueListByRepoOptions{
		State: stateArg,
	}
	if loginArg != "" {
		opts.Creator = loginArg
	}

	return issueCountBy(e, opts, func(issue *github.Issue) bool {
		return !issue.IsPullRequest()
	})
}

func issueCountBy(e aladino.Env, opts *github.IssueListByRepoOptions, predicate func(issue *github.Issue) bool) (lang.Value, error) {
	entity := e.GetTarget().GetTargetEntity()
	owner := entity.Owner
	repo := entity.Repo

	issues, _, err := e.GetGithubClient().ListIssuesByRepo(e.GetCtx(), owner, repo, opts)
	if err != nil {
		return nil, err
	}

	count := 0
	for _, issue := range issues {
		if predicate(issue) {
			count++
		}
	}

	return lang.BuildIntValue(count), nil
}
