// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	_ "github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Size() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildIntType()),
		Code:           sizeCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func sizeCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	filePatternsRegex := args[0].(*aladino.ArrayValue)

	if len(filePatternsRegex.Vals) == 0 {
		size := e.GetTarget().(*target.PullRequestTarget).PullRequest.GetAdditions() + e.GetTarget().(*target.PullRequestTarget).PullRequest.GetDeletions()
		return aladino.BuildIntValue(size), nil
	}

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	clearedPatch := make(target.Patch)
	for fp, file := range patch {
		var found bool
		for _, filePatternRegex := range filePatternsRegex.Vals {
			re, err := doublestar.Match(filePatternRegex.(*aladino.StringValue).Val, fp)
			if err != nil {
				found = false
			}
			if re {
				found = true
			}
		}

		if !found {
			clearedPatch[fp] = file
		}
	}

	size := 0
	for _, file := range clearedPatch {
		if file.Repr.Changes != nil {
			size += *file.Repr.Changes
		}
	}

	return aladino.BuildIntValue(size), nil
}
