// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"time"

	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func LastEventAt() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, aladino.BuildIntType()),
		Code:           lastEventAtCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func lastEventAtCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	updatedAt, err := e.GetTarget().GetUpdatedAt()
	if err != nil {
		return nil, err
	}

	updatedAtTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", updatedAt)
	if err != nil {
		return nil, err
	}

	return aladino.BuildIntValue(int(updatedAtTime.Unix())), nil
}
