// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func Approve() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, nil),
		Code:           approveCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

func approveCode(e aladino.Env, args []lang.Value) error {
	return reviewCode(e, []lang.Value{
		&lang.StringValue{Val: "APPROVE"},
		args[0],
	})
}
