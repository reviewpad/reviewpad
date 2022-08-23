// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"strings"

	"github.com/reviewpad/host-event-handler/handler"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func HasFileExtensions() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code:           hasFileExtensionsCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasFileExtensionsCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	argExtensions := args[0].(*aladino.ArrayValue)

	extensionSet := make(map[string]bool, len(argExtensions.Vals))
	for _, argExt := range argExtensions.Vals {
		argStringVal := argExt.(*aladino.StringValue)

		normalizedStr := strings.ToLower(argStringVal.Val)
		extensionSet[normalizedStr] = true
	}

	patch := e.GetTarget().(*target.PullRequestTarget).Patch
	for fp := range patch {
		fpExt := utils.FileExt(fp)
		normalizedExt := strings.ToLower(fpExt)

		if _, ok := extensionSet[normalizedExt]; !ok {
			return aladino.BuildFalseValue(), nil
		}
	}

	return aladino.BuildTrueValue(), nil
}
