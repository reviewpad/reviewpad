// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
)

func Fail() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code: failCode,
	}
}

func failCode(e aladino.Env, _ []aladino.Value) error {
    return fmt.Errorf("action failed")
}
