// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/stretchr/testify/assert"
)

// Use for test only
type mockUnaryOperator struct{}

// Use for test only
type mockBinaryOperator struct{}

func (op *mockUnaryOperator) getOperator() string {
	return "MOCK_UNARY_OPERATOR"
}

func (op *mockBinaryOperator) getOperator() string {
	return "MOCK_BINARY_OPERATOR"
}

func (op *mockUnaryOperator) Eval(exprVal lang.Value) lang.Value {
	return nil
}

func (op *mockBinaryOperator) Eval(lhs, rhs lang.Value) lang.Value {
	return nil
}

func TestTypeInference_WhenGivenNonExistingBuiltIn(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	expr := BuildVariable("nonBuiltIn")

	gotType, err := TypeInference(mockedEnv, expr)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInference_WhenGivenBoolConst(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil, MockBuiltIns(), nil)

	expr := BuildBoolConst(true)

	wantType := lang.BuildBoolType()

	gotType, err := TypeInference(mockedEnv, expr)

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType, "bool type is expected")
}

func TestTypesInfer_WhenGivenArrayOfExprThatContainsNonExistingBuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	exprs := []Expr{BuildVariable("nonBuiltIn")}

	gotType, err := typesinfer(mockedTypeEnv, exprs)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypesInfer_WhenGivenArrayOfExprThatContainsExistingBuiltInWithArgs(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	// returnStr is a mocked built-in that receives and returns a string value
	exprs := []Expr{BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildStringConst("hello")})}

	gotType, err := typesinfer(mockedTypeEnv, exprs)

	wantType := []lang.Type{lang.BuildStringType()}

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenUnaryOpExprIsANonExistingBuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	unaryOp := BuildUnaryOp(notOperator(), BuildVariable("nonBuiltIn"))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenUnaryOpOperatorIsNotANotOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	unaryOp := BuildUnaryOp(&mockUnaryOperator{}, BuildBoolConst(true))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "type inference failed")
}

func TestTypeInfer_WhenUnaryOpOperatorIsANotOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	gotType, err := unaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpLhsHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildVariable("nonBuiltIn"), eqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenBinaryOpRhsHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), eqOperator(), BuildVariable("nonBuiltIn"))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenBinaryOpHasEqOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), eqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasNeqOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), neqOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasGreaterEqThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), greaterEqThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasGreaterThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), greaterThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasLessEqThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), lessEqThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasLessThanOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildIntConst(1), lessThanOperator(), BuildIntConst(1))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasAndOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), andOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpHasOrOperator(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), orOperator(), BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBinaryOpOperatorIsNotAValidOp(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	binaryOp := BuildBinaryOp(BuildBoolConst(true), &mockBinaryOperator{}, BuildBoolConst(true))
	gotType, err := binaryOp.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "type inference failed")
}

func TestTypeInfer_WhenFunctionCallArgsHasTypeError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	// returnStr is a mocked built-in that receives and returns a string value
	fc := BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildVariable("nonBuiltIn")})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenFunctionCallNameHasTypeError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	fc := BuildFunctionCall(BuildVariable("nonBuiltIn"), []Expr{})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenFunctionCallHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	// returnStr is a mocked built-in that receives and returns a string value
	fc := BuildFunctionCall(BuildVariable("returnStr"), []Expr{BuildStringConst("hello")})
	gotType, err := fc.typeinfer(mockedTypeEnv)

	wantType := lang.BuildStringType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenFunctionCallHasMismatchInArgTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	// returnStr is a mocked built-in that receives and returns a string value
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
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenLambdaBodyTypeHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{},
		BuildVariable("nonBuiltIn"),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenLambdaHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{},
		BuildStringConst("hello"),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	wantType := lang.BuildFunctionType([]lang.Type{}, lang.BuildStringType())

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenLambdaHasCorrectArgumentTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	lambda := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("x"), lang.BuildStringType())},
		BuildVariable("x"),
	)
	gotType, err := lambda.typeinfer(mockedTypeEnv)

	wantType := lang.BuildFunctionType([]lang.Type{lang.BuildStringType()}, lang.BuildStringType())

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenTypedExprExprIsNotVariable(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	typedExpr := BuildTypedExpr(BuildIntConst(1), lang.BuildIntType())
	gotType, err := typedExpr.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, fmt.Sprintf("typed expression %v is not a variable", typedExpr.expr))
}

func TestTypeInfer_WhenTypedExprHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	variableName := "dummyVariable"

	typedExpr := BuildTypedExpr(BuildVariable(variableName), lang.BuildIntType())
	gotType, err := typedExpr.typeinfer(mockedTypeEnv)

	gotStoredType, ok := mockedTypeEnv[variableName]

	wantType := lang.BuildIntType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
	assert.True(t, ok)
	assert.Equal(t, wantType, gotStoredType)
}

func TestTypeInfer_WhenVariableIsNotABuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	variable := BuildVariable("nonBuiltIn")
	gotType, err := variable.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenVariableIsABuiltIn(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	// zeroConst is a mocked built-in that is a constant function of type int
	variable := BuildVariable("zeroConst")
	gotType, err := variable.typeinfer(mockedTypeEnv)

	wantType := lang.BuildFunctionType([]lang.Type{}, lang.BuildIntType())

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenStringConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	stringConst := BuildStringConst("hello")
	gotType, err := stringConst.typeinfer(mockedTypeEnv)

	wantType := lang.BuildStringType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenIntConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	intConst := BuildIntConst(1)
	gotType, err := intConst.typeinfer(mockedTypeEnv)

	wantType := lang.BuildIntType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenBoolConst(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	boolConst := BuildBoolConst(true)
	gotType, err := boolConst.typeinfer(mockedTypeEnv)

	wantType := lang.BuildBoolType()

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}

func TestTypeInfer_WhenArrayElemsTypeHasError(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	array := BuildArray([]Expr{BuildVariable("nonBuiltIn")})
	gotType, err := array.typeinfer(mockedTypeEnv)

	assert.Nil(t, gotType)
	assert.EqualError(t, err, "no type for built-in nonBuiltIn")
}

func TestTypeInfer_WhenArrayElemsTypeHasCorrectTypes(t *testing.T) {
	mockedTypeEnv := MockTypeEnv()

	array := BuildArray([]Expr{BuildIntConst(1)})
	gotType, err := array.typeinfer(mockedTypeEnv)

	wantType := lang.BuildArrayType([]lang.Type{lang.BuildIntType()})

	assert.Nil(t, err)
	assert.Equal(t, wantType, gotType)
}
