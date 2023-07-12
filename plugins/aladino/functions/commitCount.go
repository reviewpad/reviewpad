// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func CommitCount() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{}, lang.BuildIntType()),
		Code:           commitCountCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest},
	}
}

func commitCountCode(e aladino.Env, _ []lang.Value) (lang.Value, error) {
	t := e.GetTarget().(*target.PullRequestTarget)

	return lang.BuildIntValue(int(t.GetCommitCount())), nil
}
