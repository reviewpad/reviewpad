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

func Commits() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           commitsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func commitsCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)
	ghCommits, err := t.GetCommits()
	if err != nil {
		return nil, err
	}

	commitMessages := make([]lang.Value, len(ghCommits))
	for i, ghCommit := range ghCommits {
		commitMessages[i] = lang.BuildStringValue(ghCommit.Message)
	}

	return lang.BuildArrayValue(commitMessages), nil
}
