// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"strings"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Join() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType()), aladino.BuildStringType()}, aladino.BuildStringType()),
		Code:           joinCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func joinCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	elementsArg := args[0].(*lang.ArrayValue)
	if len(elementsArg.Vals) == 0 {
		return lang.BuildStringValue(""), nil
	}
	var clearVals []string
	separatorArg := args[1].(*lang.StringValue)
	for _, val := range elementsArg.Vals {
		switch v := val.(type) {
		case *lang.StringValue:
			clearVals = append(clearVals, v.Val)
		default:
			return nil, fmt.Errorf("join: invalid element of kind %v", v.Kind())
		}
	}

	return lang.BuildStringValue(strings.Join(clearVals, separatorArg.Val)), nil
}
