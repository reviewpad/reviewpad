// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	consts "github.com/reviewpad/reviewpad/v3/plugins/aladino/consts"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func ErrorMsg() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: errorCode,
	}
}

func errorCode(e aladino.Env, args []aladino.Value) error {
	body := args[0].(*aladino.StringValue).Val

	comments := e.GetComments()
	errors, ok := comments[consts.ERROR_LEVEL]
	if !ok {
		comments[consts.ERROR_LEVEL] = []string{body}
	} else {
		comments[consts.ERROR_LEVEL] = append(errors, body)
	}

	return nil
}
