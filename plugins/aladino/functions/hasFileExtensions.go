// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
)

func HasFileExtensions() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code:           hasFileExtensionsCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func hasFileExtensionsCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	argExtensions := args[0].(*lang.ArrayValue)

	extensionSet := make(map[string]bool, len(argExtensions.Vals))
	for _, argExt := range argExtensions.Vals {
		argStringVal := argExt.(*lang.StringValue)

		normalizedStr := strings.ToLower(argStringVal.Val)
		extensionSet[normalizedStr] = true
	}

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	for fp := range patch {
		fpExt := utils.FileExt(fp)
		normalizedExt := strings.ToLower(fpExt)

		if _, ok := extensionSet[normalizedExt]; !ok {
			return lang.BuildFalseValue(), nil
		}
	}

	return lang.BuildTrueValue(), nil
}
