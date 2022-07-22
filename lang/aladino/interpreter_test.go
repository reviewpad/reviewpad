// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestEvalExpr_WhenParseFails(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	gotVal, err := aladino.EvalExpr(mockedEnv, "", "1 ==")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "parse error: failed to build AST on input 1 ==")
}

func TestEvalExpr_WhenTypeInferenceFails(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	gotVal, err := aladino.EvalExpr(mockedEnv, "", "1 == \"a\"")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "type inference failed")
}

func TestEvalExpr_WhenExprIsNotBoolType(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	gotVal, err := aladino.EvalExpr(mockedEnv, "", "1")

	assert.False(t, gotVal)
	assert.EqualError(t, err, "expression 1 is not a condition")
}

func TestEvalExpr(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	gotVal, err := aladino.EvalExpr(mockedEnv, "", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestEvalExpr_OnInterpreter(t *testing.T) {
	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	mockedInterpreter := &aladino.Interpreter{
		Env: mockedEnv,
	}

	gotVal, err := mockedInterpreter.EvalExpr("", "1 == 1")

	assert.Nil(t, err)
	assert.True(t, gotVal)
}
