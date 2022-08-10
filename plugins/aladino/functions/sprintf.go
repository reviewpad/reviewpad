// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Sprintf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildStringType()),
		Code: sprintfCode,
	}
}

func sprintfCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	var clearVals []interface{}
	format := args[0].(*aladino.StringValue).Val
	vals := args[1].(*aladino.ArrayValue).Vals

	for _, val := range vals {
		switch val.(type) {
		case *aladino.StringValue:
			clearVals = append(clearVals, val.(*aladino.StringValue).Val)
		}
	}

	return aladino.BuildStringValue(fmt.Sprintf(format, clearVals...)), nil
}
