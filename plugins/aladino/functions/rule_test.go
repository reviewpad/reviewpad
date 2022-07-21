// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var rule = plugins_aladino.PluginBuiltIns().Functions["rule"].Code

func TestRule_WhenRuleIsAbsent(t *testing.T) {
	ruleName := "is-absent"
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, fmt.Sprintf("$rule: no rule with name %v in state %+q", ruleName, mockedEnv.GetRegisterMap()))
}

func TestRule_WhenRuleIsInvalid(t *testing.T) {
	ruleName := "is-invalid"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == \"a\"")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestRule_WhenRuleIsTrue(t *testing.T) {
	ruleName := "tautology"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == 1")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestRule_WhenRuleIsFalse(t *testing.T) {
	ruleName := "is-false-premise"
	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedEnv.GetRegisterMap()[internalRuleName] = aladino.BuildStringValue("1 == 2")

	args := []aladino.Value{aladino.BuildStringValue(ruleName)}
	gotVal, err := rule(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
