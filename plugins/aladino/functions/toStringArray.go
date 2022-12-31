// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"encoding/json"
	"fmt"

	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func ToStringArray() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
		Code:           toStringArray,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func toStringArray(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	str := args[0].(*aladino.StringValue).Val
	arr := []string{}
	elements := []aladino.Value{}

	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return nil, fmt.Errorf(`error converting "%s" to string array: %s`, str, err.Error())
	}

	for _, value := range arr {
		elements = append(elements, aladino.BuildStringValue(value))
	}

	return aladino.BuildArrayValue(elements), nil
}
