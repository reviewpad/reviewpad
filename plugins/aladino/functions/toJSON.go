// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/ohler55/ojg/oj"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func ToJSON() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildJSONType()),
		Code:           toJSONCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func toJSONCode(e aladino.Env, args []lang.Value) (lang.Value, error) {
	jsonString := args[0].(*lang.StringValue).Val

	val, err := oj.ParseString(jsonString)
	if err != nil {
		return nil, err
	}

	return lang.BuildJSONValue(val), nil
}
