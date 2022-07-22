// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestProcessRule(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	ruleName := "rule_name"
	spec := "1 == 1"
	err = mockedInterpreter.ProcessRule(ruleName, spec)

	internalRuleName := fmt.Sprintf("@rule:%v", ruleName)
	gotVal := mockedEnv.GetRegisterMap()[internalRuleName]

	wantVal := aladino.BuildStringValue(spec)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
