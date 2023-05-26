// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var disableActions = plugins_aladino.PluginBuiltIns().Actions["disableActions"].Code

func TestDisableActions(t *testing.T) {
	builtInName := "emptyAction"

	builtIns := &aladino.BuiltIns{
		Actions: map[string]*aladino.BuiltInAction{
			builtInName: {
				Type: lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildArrayOfType(lang.BuildStringType())),
				Code: func(e aladino.Env, args []lang.Value) error {
					return nil
				},
				Disabled: false,
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, builtIns, nil)

	args := []lang.Value{lang.BuildArrayValue([]lang.Value{lang.BuildStringValue(builtInName)})}
	err := disableActions(mockedEnv, args)

	assert.Nil(t, err)
	assert.True(t, mockedEnv.GetBuiltIns().Actions[builtInName].Disabled)
}
