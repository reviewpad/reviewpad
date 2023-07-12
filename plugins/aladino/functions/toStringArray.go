// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"encoding/json"
	"fmt"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ToStringArray() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
		Code:           toStringArray,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func toStringArray(e aladino.Env, args []lang.Value) (lang.Value, error) {
	str := args[0].(*lang.StringValue).Val
	arr := []string{}
	elements := []lang.Value{}

	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return nil, fmt.Errorf(`error converting "%s" to string array: %s`, str, err.Error())
	}

	for _, value := range arr {
		elements = append(elements, lang.BuildStringValue(value))
	}

	return lang.BuildArrayValue(elements), nil
}
