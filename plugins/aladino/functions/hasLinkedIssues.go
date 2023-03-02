// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func HasLinkedIssues() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildBoolType()),
		Code:           hasLinkedIssuesCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasLinkedIssuesCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	closingIssuesCount, err := pullRequest.GetLinkedIssuesCount()
	if err != nil {
		return nil, err
	}

	return aladino.BuildBoolValue(closingIssuesCount > 0), nil
}
