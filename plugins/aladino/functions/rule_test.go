// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var rule = plugins_aladino.PluginBuiltIns().Functions["rule"].Code

func TestRule_WhenRuleIsAbsent(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ruleName := "is-absent"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, fmt.Sprintf("$rule: no rule with name %v in state %+q", ruleName, mockedEnv.GetRegisterMap()))
}

func TestRule_WhenRuleIsInvalid(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ruleName := "is-invalid"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == \"a\"")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestRule_WhenRuleIsTrue(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ruleName := "tautology"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == 1")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestRule_WhenRuleIsFalse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ruleName := "is-false-premise"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == 2")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
