// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	consts "github.com/reviewpad/reviewpad/v3/plugins/aladino/consts"
)

func Info() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code: infoCode,
	}
}

func infoCode(e aladino.Env, args []aladino.Value) error {
	body := args[0].(*aladino.StringValue).Val

	comments := e.GetComments()
	infos, ok := comments[consts.INFO_LEVEL]
	if !ok {
		comments[consts.INFO_LEVEL] = []string{body}
	} else {
		comments[consts.INFO_LEVEL] = append(infos, body)
	}

	return nil
}
