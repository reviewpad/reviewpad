// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

var (
	// this regexp matches the golang formatting verbs such as %v, %s...
	verbsRegExp = regexp.MustCompile(`%(\#|\+|\-| |0)?(\[\d+\])?(([1-9])\.([1-9])|([1-9])|([1-9])\.|\.([1-9]))?(\w{1,9})`)
)

func Sprintf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildDynamicArrayType()}, aladino.BuildStringType()),
		Code:           sprintfCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func sprintfCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	var clearVals []interface{}
	format := args[0].(*aladino.StringValue).Val
	vals := args[1].(*aladino.ArrayValue).Vals

	verbs := verbsRegExp.FindAllString(format, -1)

	if len(vals) > len(verbs) {
		vals = vals[0:len(verbs)]
	}

	for _, val := range vals {
		switch v := val.(type) {
		case *aladino.StringValue:
			clearVals = append(clearVals, v.Val)
		case *aladino.IntValue:
			clearVals = append(clearVals, v.Val)
		case *aladino.BoolValue:
			clearVals = append(clearVals, v.Val)
		}
	}

	return aladino.BuildStringValue(fmt.Sprintf(format, clearVals...)), nil
}
