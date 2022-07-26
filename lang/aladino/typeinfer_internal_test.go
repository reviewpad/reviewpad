// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nilUnaryOp struct{}

type nilBinaryOp struct{}

func (op *nilUnaryOp) getOperator() string {
	return "NIL_UNARY_OP"
}

func (op *nilBinaryOp) getOperator() string {
	return "NIL_BINARY_OP"
}

func (op *nilUnaryOp) Eval(exprVal Value) Value {
	return nil
}

func (op *nilBinaryOp) Eval(lhs, rhs Value) Value {
	return nil
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
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypesInfer_WhenGivenArrayOfExprThatContainsNonExistingBuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	exprs := []Expr{BuildVariable("nonBuiltIn")}

	gotType, err := typesinfer(mockedTypeEnv, exprs)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
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
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenUnaryOpOperatorIsNotANotOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	unaryOp := BuildUnaryOp(&nilUnaryOp{}, BuildBoolConst(true))
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

func TestTypeInfer_WhenBinaryOpLhsHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildVariable("nonBuiltIn"), eqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenBinaryOpRhsHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), eqOperator(), BuildVariable("nonBuiltIn"))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenBinaryOpHasEqOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), eqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasNeqOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), neqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasGreaterEqThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), greaterEqThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasGreqaterThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), greaterThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasLessEqThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), lessEqThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasLessThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), lessThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasAndOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), andOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasOrOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), orOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpOperatorIsNotAValidOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), &nilBinaryOp{}, BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "type inference failed")
}

func TestTypeInfer_WhenFunctionCallArgsHasTypeError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	fc := BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildVariable("nonBuiltIn")})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenFunctionCallNameHasTypeError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	fc := BuildFunctionCall(BuildVariable("nonBuiltIn"), []Expr{})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenFunctionCallHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	fc := BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildStringConst("hello")})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	wantType := BuildStringType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenFunctionCallHasMismatchInArgTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	fc := BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildIntConst(1)})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "type inference failed: mismatch in arg types on returnStr")
}

func TestTypeInfer_WhenLambdaParamTypeHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{BuildVariable("nonBuiltIn")},
		BuildIntConst(1),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenLambdaBodyTypeHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{BuildIntConst(1)},
		BuildVariable("nonBuiltIn"),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenLambdaHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{BuildVariable("returnStr")},
		BuildStringConst("hello"),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	wantType := BuildFunctionType([]Type{BuildFunctionType([]Type{BuildStringType()}, BuildStringType())}, BuildStringType())

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenTypedExprExprIsNotVariable(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	gotType, err := typedExpr.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, fmt.Sprintf("typed expression %v is not a variable", typedExpr.expr))
}

func TestTypeInfer_WhenTypedExprHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	typedExpr := BuildTypedExpr(BuildVariable("zeroConst"), BuildIntType())
	gotType, err := typedExpr.typeinfer(mockedTypeEnv)

	wantType := BuildIntType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenVariableIsNotABuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	variable := BuildVariable("nonBuiltIn")
	gotType, err := variable.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenVariableIsABuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	variable := BuildVariable("zeroConst")
	gotType, err := variable.typeinfer(mockedTypeEnv)

	wantType := BuildFunctionType([]Type{}, BuildIntType())

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenInferingStringConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	stringConst := BuildStringConst("hello")
	gotType, err := stringConst.typeinfer(mockedTypeEnv)

	wantType := BuildStringType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenInferingIntConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	intConst := BuildIntConst(1)
	gotType, err := intConst.typeinfer(mockedTypeEnv)

	wantType := BuildIntType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenInferingBoolConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	boolConst := BuildBoolConst(true)
	gotType, err := boolConst.typeinfer(mockedTypeEnv)

	wantType := BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenArrayElemsTypeHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	array := BuildArray([]Expr{BuildVariable("nonBuiltIn")})
	gotType, err := array.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn. Please check if the mode in the reviewpad.yml file supports it")
}

func TestTypeInfer_WhenArrayElemsTypeHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	array := BuildArray([]Expr{BuildIntConst(1)})
	gotType, err := array.typeinfer(mockedTypeEnv)

	wantType := BuildArrayType([]Type{BuildIntType()})

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}
