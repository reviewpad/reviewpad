// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"time"

	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func CreatedAt() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code:           createdAtCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func createdAtCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	createdAt, err := e.GetTarget().GetCreatedAt()
	if err != nil {
		return nil, err
	}

	createdAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", createdAt)
	if err != nil {
		return nil, err
	}
	return aladino.BuildIntValue(int(createdAtTime.Unix())), nil
}
