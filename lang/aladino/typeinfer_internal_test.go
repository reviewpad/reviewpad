// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nilOp struct {}

func (op *nilOp) getOperator() string {
    return "NIL_OP"
}

func (op *nilOp) Eval(exprVal Value) Value {
	return BuildBoolValue(!exprVal.(*BoolValue).Val)
}

func TestTypeInference_WhenGivenBoolConst(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	expr := BuildBoolConst(true)

	wantType := BuildBoolType()

	gotType, err := TypeInference(mockedEnv, expr)

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType, "bool type is expected")
}

func TestTypeInference_WhenGivenNonExistingBuiltIn(t *testing.T) {
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		assert.FailNow(t, "MockDefaultEnv returned unexpected error: %v", err)
	}

	expr := BuildVariable("nonBuiltIn")

	gotType, err := TypeInference(mockedEnv, expr)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it.")
}

func TestTypesInfer_WhenGivenArrayOfExprThatContainsNonExistingBuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	exprs := []Expr{BuildVariable("nonBuiltIn")}

	gotType, err := typesinfer(mockedTypeEnv, exprs)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it.")
}

func TestTypesInfer_WhenGivenArrayOfExprThatContainsExistingBuiltInWithArgs(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	exprs := []Expr{BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildStringConst("hello")})}

	gotType, err := typesinfer(mockedTypeEnv, exprs)

	wantType := []Type{BuildStringType()}

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenUnaryOpExprIsANonExistingBuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	unaryOp := BuildUnaryOp(notOperator(), BuildVariable("nonBuiltIn"))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it.")
}

func TestTypeInfer_WhenUnaryOpOperatorIsNotANotOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()
    
	unaryOp := BuildUnaryOp(&nilOp{}, BuildBoolConst(true))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "type inference failed")
}

func TestTypeInfer_WhenUnaryOpOperatorIsANotOp(t *testing.T) {
    mockedTypeEnv := MockTypeEnv()
    
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

    wantType := BuildBoolType()

	assert.Nil(t, err)
    assert.Equal(t, wantType, gotType)
}
