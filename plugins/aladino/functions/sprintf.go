// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"regexp"

	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

var (
	// this regexp matches the golang formatting verbs such as %v, %s...
	verbsRegExp = regexp.MustCompile(`%(\#|\+|\-| |0)?(\[\d+\])?(([1-9])\.([1-9])|([1-9])|([1-9])\.|\.([1-9]))?(\w{1,9})`)
)

func Sprintf() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildDynamicArrayType()}, lang.BuildStringType()),
		Code:           sprintfCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
	}
}

func sprintfCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	var clearVals []interface{}
	format := args[0].(*lang.StringValue).Val
	vals := args[1].(*lang.ArrayValue).Vals

	verbs := verbsRegExp.FindAllString(format, -1)

	if len(vals) > len(verbs) {
		vals = vals[0:len(verbs)]
	}

	for _, val := range vals {
		switch v := val.(type) {
		case *lang.StringValue:
			clearVals = append(clearVals, v.Val)
		case *lang.IntValue:
			clearVals = append(clearVals, v.Val)
		case *lang.BoolValue:
			clearVals = append(clearVals, v.Val)
		}
	}

	return lang.BuildStringValue(fmt.Sprintf(format, clearVals...)), nil
}
