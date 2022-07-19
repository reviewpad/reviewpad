// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"log"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func Fail() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: failCode,
	}
}

func failCode(e aladino.Env, args []aladino.Value) error {
	failMessage := args[0].(*aladino.StringValue).Val

	log.Fatalf(failMessage)

	return nil
}
