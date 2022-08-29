// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Commits() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           commitsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func commitsCode(e aladino.Env, _ []aladino.Value) (aladino.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)
	ghCommits, err := t.GetCommits()
	if err != nil {
		return nil, err
	}

	commitMessages := make([]aladino.Value, len(ghCommits))
	for i, ghCommit := range ghCommits {
		commitMessages[i] = aladino.BuildStringValue(ghCommit.Message)
	}

	return aladino.BuildArrayValue(commitMessages), nil
}
