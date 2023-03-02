// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Approve() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           approveCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func approveCode(e aladino.Env, args []aladino.Value) error {
	return reviewCode(e, []aladino.Value{
		&aladino.StringValue{Val: "APPROVE"},
		args[0],
	})
}
