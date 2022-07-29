// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var disableActions = plugins_aladino.PluginBuiltIns().Actions["disableActions"].Code

func TestDisableActions(t *testing.T) {
	builtInName := "emptyAction"

	builtIns := &aladino.BuiltIns{
		Actions: map[string]*aladino.BuiltInAction{
			builtInName: {
				Type: aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildArrayOfType(aladino.BuildStringType())),
				Code: func(e aladino.Env, args []aladino.Value) error {
					return fmt.Errorf("%v is called", builtInName)
				},
				Disabled: false,
			},
		},
	}
	mockedEnv, err := aladino.MockDefaultEnvWithBuiltIns(nil, nil, builtIns)
	if err != nil {
		log.Fatalf("mockDefaultEnvWithBuiltIns failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue(builtInName)})}
	err = disableActions(mockedEnv, args)

	assert.Nil(t, err)
	assert.True(t, mockedEnv.GetBuiltIns().Actions[builtInName].Disabled)
}
