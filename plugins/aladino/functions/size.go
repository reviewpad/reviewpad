// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Size() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildIntType()),
		Code:           sizeCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func sizeCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	filePatternsRegex := args[0].(*lang.ArrayValue)

	if len(filePatternsRegex.Vals) == 0 {
		size := e.GetTarget().(*target.PullRequestTarget).PullRequest.AdditionsCount + e.GetTarget().(*target.PullRequestTarget).PullRequest.DeletionsCount
		return lang.BuildIntValue(int(size)), nil
	}

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	clearedPatch := make(target.Patch)
	for fp, file := range patch {
		var found bool
		for _, filePatternRegex := range filePatternsRegex.Vals {
			re, err := doublestar.Match(filePatternRegex.(*lang.StringValue).Val, fp)
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
		size += int(file.Repr.ChangesCount)
	}

	return lang.BuildIntValue(size), nil
}
