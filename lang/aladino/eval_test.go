// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestEval_OnUnaryOp_WhenExprEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	unaryOp, err := aladino.Parse("!$nonBuiltIn")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := unaryOp.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnUnaryOp(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	unaryOp, err := aladino.Parse("!true")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := unaryOp.Eval(mockedEnv)

	wantVal := lang.BuildFalseValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnBinaryOp_WhenLeftOperandEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	binaryOp, err := aladino.Parse("$nonBuiltIn() == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := binaryOp.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnBinaryOp_WhenRightOperandEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	binaryOp, err := aladino.Parse("1 == $nonBuiltIn()")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := binaryOp.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnBinaryOp_WhenDiffKinds(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	binaryOp, err := aladino.Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := binaryOp.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: left and right operand have different kinds")
}

func TestEval_OnBinaryOp_WhenTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	binaryOp, err := aladino.Parse("1 == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := binaryOp.Eval(mockedEnv)

	wantVal := lang.BuildTrueValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnBinaryOp_WhenFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	binaryOp, err := aladino.Parse("1 == 2")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := binaryOp.Eval(mockedEnv)

	wantVal := lang.BuildFalseValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnVariable_WhenVariableIsRegistered(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	variable, err := aladino.Parse("$size")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	variableName := "size"
	mockedEnv.GetRegisterMap()[variableName] = lang.BuildIntValue(0)

	gotVal, err := variable.Eval(mockedEnv)

	wantVal := lang.BuildIntValue(0)

	// clean up
	delete(mockedEnv.GetRegisterMap(), variableName)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnVariable_WhenVariableIsNotABuiltIn(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	variable, err := aladino.Parse("$nonBuiltIn")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := variable.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnVariable_WhenVariableIsABuiltIn(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	variable, err := aladino.Parse("$zeroConst")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := variable.Eval(mockedEnv)

	wantVal := lang.BuildIntValue(0)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnBoolConst(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	boolConst, err := aladino.Parse("true")
	if err != nil {
		assert.FailNow(t, "parse fa", err)
	}

	gotVal, err := boolConst.Eval(mockedEnv)

	wantVal := lang.BuildTrueValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnStringConst(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	str := "test string"
	strConst, err := aladino.Parse("\"" + str + "\"")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := strConst.Eval(mockedEnv)

	wantVal := lang.BuildStringValue(str)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnIntConst(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	intConst, err := aladino.Parse("1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := intConst.Eval(mockedEnv)

	wantVal := lang.BuildIntValue(1)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnFunctionCall_WhenArgEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	fc, err := aladino.Parse("$addLabel($nonBuiltIn)")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := fc.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnFunctionCall_WhenFunctionIsNotABuiltIn(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	fc, err := aladino.Parse("$nonBuiltIn()")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := fc.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnFunctionCall_WhenFunctionIsABuiltIn(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	fc, err := aladino.Parse("$returnStr(\"hello\")")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := fc.Eval(mockedEnv)

	wantVal := lang.BuildStringValue("hello")

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnRecursiveFunctionCall_WhenFunctionIsABuiltIn(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	fc, err := aladino.Parse("$returnStr($returnStr(\"hello\"))")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := fc.Eval(mockedEnv)

	wantVal := lang.BuildStringValue("hello")

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnLambda_WhenLambdaBodyEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	lambdaBody, err := aladino.Parse("1 == $nonBuiltIn")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	lambda := aladino.BuildLambda([]aladino.Expr{}, lambdaBody)

	gotFn, err := lambda.Eval(mockedEnv)

	gotVal := gotFn.(*lang.FunctionValue).Fn([]lang.Value{})

	assert.Nil(t, err)
	assert.Nil(t, gotVal)
}

func TestEval_OnLambda(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	variable, err := aladino.Parse("$size")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}
	lambdaParam := aladino.BuildTypedExpr(variable, aladino.BuildStringType())
	lambdaBody, err := aladino.Parse("1 == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}
	lambda := aladino.BuildLambda([]aladino.Expr{lambdaParam}, lambdaBody)

	gotFn, err := lambda.Eval(mockedEnv)

	gotVal := gotFn.(*lang.FunctionValue).Fn([]lang.Value{lang.BuildIntValue(0)})

	wantVal := lang.BuildTrueValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnTypedExpr(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	typedExpr := aladino.BuildTypedExpr(expr, aladino.BuildBoolType())

	gotVal, err := typedExpr.Eval(mockedEnv)

	wantVal := lang.BuildTrueValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnArray_WhenElemEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	array, err := aladino.Parse("[$nonBuiltIn()]")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := array.Eval(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: failure on nonBuiltIn")
}

func TestEval_OnArray(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	array, err := aladino.Parse("[\"a\"]")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := array.Eval(mockedEnv)

	wantVal := lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("a")})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEval_WhenExprEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := aladino.Eval(mockedEnv, expr)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, "eval: left and right operand have different kinds")
}

func TestEval(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := aladino.Eval(mockedEnv, expr)

	wantVal := lang.BuildTrueValue()

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestEvalCondition_WhenExprEvalFails(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == \"a\"")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := aladino.EvalCondition(mockedEnv, expr)

	assert.False(t, gotVal)
	assert.EqualError(t, err, "eval: left and right operand have different kinds")
}

func TestEvalCondition_WhenConditionIsTrue(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == 1")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := aladino.EvalCondition(mockedEnv, expr)

	assert.Nil(t, err)
	assert.True(t, gotVal)
}

func TestEvalCondition_WhenConditionIsFalse(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	expr, err := aladino.Parse("1 == 2")
	if err != nil {
		assert.FailNow(t, "parse failed", err)
	}

	gotVal, err := aladino.EvalCondition(mockedEnv, expr)

	assert.Nil(t, err)
	assert.False(t, gotVal)
}

func TestEval_OnNotOp(t *testing.T) {
	notOp := &aladino.NotOp{}
	gotVal := notOp.Eval(lang.BuildTrueValue())

	wantVal := lang.BuildFalseValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnEqOp(t *testing.T) {
	eqOp := &aladino.EqOp{}
	gotVal := eqOp.Eval(lang.BuildIntValue(1), lang.BuildIntValue(1))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnNeqOp(t *testing.T) {
	neqOp := &aladino.NeqOp{}
	gotVal := neqOp.Eval(lang.BuildIntValue(1), lang.BuildIntValue(2))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnAndOp(t *testing.T) {
	andOp := &aladino.AndOp{}
	gotVal := andOp.Eval(lang.BuildTrueValue(), lang.BuildTrueValue())

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnOrOp(t *testing.T) {
	orOp := &aladino.OrOp{}
	gotVal := orOp.Eval(lang.BuildTrueValue(), lang.BuildTrueValue())

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnLessThanOp(t *testing.T) {
	lessThanOp := &aladino.LessThanOp{}
	gotVal := lessThanOp.Eval(lang.BuildIntValue(1), lang.BuildIntValue(2))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnLessEqThanOp(t *testing.T) {
	lessEqThanOp := &aladino.LessEqThanOp{}
	gotVal := lessEqThanOp.Eval(lang.BuildIntValue(1), lang.BuildIntValue(2))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnGreaterThanOp(t *testing.T) {
	greaterThanOp := &aladino.GreaterThanOp{}
	gotVal := greaterThanOp.Eval(lang.BuildIntValue(3), lang.BuildIntValue(2))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}

func TestEval_OnGreaterEqThanOp(t *testing.T) {
	greaterEqThanOp := &aladino.GreaterEqThanOp{}
	gotVal := greaterEqThanOp.Eval(lang.BuildIntValue(3), lang.BuildIntValue(2))

	wantVal := lang.BuildTrueValue()

	assert.Equal(t, wantVal, gotVal)
}
