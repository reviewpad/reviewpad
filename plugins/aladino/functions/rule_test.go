// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var rule = plugins_aladino.PluginBuiltIns().Functions["rule"].Code

func TestRule_WhenRuleIsAbsent(t *testing.T) {
	ruleName := "is-absent"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []lang.Value{lang.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, fmt.Sprintf("$rule: no rule with name %v in state %+q", ruleName, mockedEnv.GetRegisterMap()))
}

func TestRule_WhenRuleIsInvalid(t *testing.T) {
	ruleName := "is-invalid"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	mockedEnv.GetRegisterMap()[internalRuleName] = lang.BuildStringValue("1 == \"a\"")

	args := []lang.Value{lang.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestRule_WhenRuleIsTrue(t *testing.T) {
	ruleName := "tautology"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	mockedEnv.GetRegisterMap()[internalRuleName] = lang.BuildStringValue("1 == 1")

	args := []lang.Value{lang.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := lang.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestRule_WhenRuleIsFalse(t *testing.T) {
	ruleName := "is-false-premise"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	mockedEnv.GetRegisterMap()[internalRuleName] = lang.BuildStringValue("1 == 2")

	args := []lang.Value{lang.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := lang.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
