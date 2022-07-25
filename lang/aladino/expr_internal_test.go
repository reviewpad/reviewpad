// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEqOperator_ExpectEqOperator(t *testing.T) {
	wantVal := &EqOp{}
	gotVal := eqOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestNeqOperator_ExpectNeqOperator(t *testing.T) {
	wantVal := &NeqOp{}
	gotVal := neqOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestAndOperator_ExpectAndOperator(t *testing.T) {
	wantVal := &AndOp{}
	gotVal := andOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestOrOperator_ExpectOrOperator(t *testing.T) {
	wantVal := &OrOp{}
	gotVal := orOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestLessThanOperator_ExpectLessThanOperator(t *testing.T) {
	wantVal := &LessThanOp{}
	gotVal := lessThanOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestLessEqThanOperator_ExpectLessEqOperator(t *testing.T) {
	wantVal := &LessEqThanOp{}
	gotVal := lessEqThanOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGreaterThanOperator_ExpectGreaterThanOperator(t *testing.T) {
	wantVal := &GreaterThanOp{}
	gotVal := greaterThanOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGreaterEqThanOperator_ExpectGreaterEqThanOperator(t *testing.T) {
	wantVal := &GreaterEqThanOp{}
	gotVal := greaterEqThanOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectNotOp(t *testing.T) {
	wantVal := NOT_OP
	gotVal := notOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectEqOp(t *testing.T) {
	wantVal := EQ_OP
	gotVal := eqOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectNeqOp(t *testing.T) {
	wantVal := NEQ_OP
	gotVal := neqOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectAndOp(t *testing.T) {
	wantVal := AND_OP
	gotVal := andOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectOrOp(t *testing.T) {
	wantVal := OR_OP
	gotVal := orOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectLessThanOp(t *testing.T) {
	wantVal := LESS_THAN_OP
	gotVal := lessThanOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectLessEqThanOp(t *testing.T) {
	wantVal := LESS_EQ_THAN_OP
	gotVal := lessEqThanOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectGreaterThanOp(t *testing.T) {
	wantVal := GREATER_THAN_OP
	gotVal := greaterThanOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_ExpectGreaterEqThanOp(t *testing.T) {
	wantVal := GREATER_EQ_THAN_OP
	gotVal := greaterEqThanOperator().getOperator()

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolConst_ExpectBoolConstWithTrueValue(t *testing.T) {
	wantVal := &BoolConst{true}
	gotVal := BuildBoolConst(true)

	assert.Equal(t, wantVal, gotVal)
}

func TestBoolConstKind_ExpectBoolConstKind(t *testing.T) {
	wantVal := BOOL_CONST
	gotVal := BuildBoolConst(true).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestBoolConstEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildIntConst(0)

	assert.False(t, boolConst.equals(otherConst))
}

func TestBoolConstEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildBoolConst(false)

	assert.False(t, boolConst.equals(otherConst))
}

func TestBoolConstEquals_ExpectTrue(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildBoolConst(true)

	assert.True(t, boolConst.equals(otherConst))
}

func TestBuildStringConst_ExpectStringConstWithLoremIpsumValue(t *testing.T) {
	wantVal := &StringConst{"LoremIpsum"}
	gotVal := BuildStringConst("LoremIpsum")

	assert.Equal(t, wantVal, gotVal)
}

func TestStringConstKind_ExpectStringConstKind(t *testing.T) {
	wantVal := STRING_CONST
	gotVal := BuildStringConst("Lorem Ipsum").Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestStringConstEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum")
	otherConst := BuildIntConst(0)

	assert.False(t, stringConst.equals(otherConst))
}

func TestStringConstEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum #1")
	otherConst := BuildStringConst("Lorem Ipsum #2")

	assert.False(t, stringConst.equals(otherConst))
}

func TestStringConstEquals_ExpectTrue(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum")
	otherConst := BuildStringConst("Lorem Ipsum")

	assert.True(t, stringConst.equals(otherConst))
}

func TestBuildIntConst_ExpectIntConstWithZeroValue(t *testing.T) {
	wantVal := &IntConst{0}
	gotVal := BuildIntConst(0)

	assert.Equal(t, wantVal, gotVal)
}

func TestIntConstKind_ExpectIntConstKind(t *testing.T) {
	wantVal := INT_CONST
	gotVal := BuildIntConst(0).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestIntConstEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildStringConst("Lorem Ipsum")

	assert.False(t, intConst.equals(otherConst))
}

func TestIntConstEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildIntConst(1)

	assert.False(t, intConst.equals(otherConst))
}

func TestIntConstEquals_ExpectTrue(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildIntConst(0)

	assert.True(t, intConst.equals(otherConst))
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsYear_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 year ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.AddDate(-timeValue, 0, 0).Unix()),
	}

	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsMonth_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 month ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.AddDate(0, -timeValue, 0).Unix()),
	}
	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsDay_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 day ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.AddDate(0, 0, -timeValue).Unix()),
	}

	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsWeek_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 week ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.Add(-(time.Hour * 24 * 7) * time.Duration(timeValue)).Unix()),
	}

	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsHour_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 hour ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.Add(-time.Hour * time.Duration(timeValue)).Unix()),
	}

	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsMinute_ExpectIntConstWithAValue(t *testing.T) {
	now := time.Now()

	val := "1 minute ago"

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, _ := strconv.Atoi(timeValueRegex.FindString(val))

	wantVal := &IntConst{
		value: int(now.Add(-time.Minute * time.Duration(timeValue)).Unix()),
	}

	gotVal := BuildRelativeTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildTimeConst_ExpectIntConstWithAValue(t *testing.T) {
	val := "2019-10-12T07:50:52"

	wantVal := &IntConst{
		value: int(time.Date(2019, time.Month(10), 12, 7, 50, 52, 0, time.UTC).Unix()),
	}

	gotVal := BuildTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildVariable_ExpectVariableWithTestValue(t *testing.T) {
	wantVal := &Variable{"Test"}
	gotVal := BuildVariable("Test")

	assert.Equal(t, wantVal, gotVal)
}

func TestVariableKind_ExpectVariableKind(t *testing.T) {
	wantVal := VARIABLE_CONST
	gotVal := BuildVariable("test").Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestVariableEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	variable := BuildVariable("test")
	otherVal := BuildIntConst(0)

	assert.False(t, variable.equals(otherVal))
}

func TestVariableEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	variable := BuildVariable("test#1")
	otherVal := BuildVariable("test#2")

	assert.False(t, variable.equals(otherVal))
}

func TestVariableEquals_ExpectTrue(t *testing.T) {
	variable := BuildVariable("test")
	otherVal := BuildVariable("test")

	assert.True(t, variable.equals(otherVal))
}

func TestBuildUnaryOp_ExpectUnaryOpWithNotOpAndBoolConstWithTrueValue(t *testing.T) {
	wantVal := &UnaryOp{&NotOp{}, &BoolConst{true}}
	gotVal := BuildUnaryOp(notOperator(), BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildNotOp_ExpectNotOpWithBoolContWithTrueValue(t *testing.T) {
	wantVal := &UnaryOp{&NotOp{}, &BoolConst{true}}
	gotVal := BuildNotOp(BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestUnaryOpKind_WhenUnaryOpHasNotOpWithBoolContWithTrueValue_ExpectUnaryOpKind(t *testing.T) {
	wantVal := UNARY_OP_CONST
	gotVal := BuildUnaryOp(notOperator(), BuildBoolConst(true)).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestUnaryOpEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildIntConst(0)

	assert.False(t, unaryOp.equals(otherVal))
}

func TestUnaryOpEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildUnaryOp(notOperator(), BuildBoolConst(false))

	assert.False(t, unaryOp.equals(otherVal))
}

func TestUnaryOpEquals_ExpectTrue(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildUnaryOp(notOperator(), BuildBoolConst(true))

	assert.True(t, unaryOp.equals(otherVal))
}

func TestBuildBinaryOp_ExpectEqOpBetweenIntConstWithOneValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &EqOp{}, &IntConst{1}}
	gotVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildAndOp_ExpectAndOpBetweenBoolConstsWithTrueValue(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &AndOp{}, &BoolConst{true}}
	gotVal := BuildAndOp(BuildBoolConst(true), BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildOrOp_ExpectOrOpBetweenBoolConstsWithTrueValue(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &OrOp{}, &BoolConst{true}}
	gotVal := BuildOrOp(BuildBoolConst(true), BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildEqOp_ExpectEqOpBetweenBoolConstsWithTrueValue(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &EqOp{}, &BoolConst{true}}
	gotVal := BuildEqOp(BuildBoolConst(true), BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildNeqOp_ExpectNeqOpBetweenBoolConstsWithTrueValue(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &NeqOp{}, &BoolConst{true}}
	gotVal := BuildNeqOp(BuildBoolConst(true), BuildBoolConst(true))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildLessThanOp_ExpectLessThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessThanOp{}, &IntConst{2}}
	gotVal := BuildLessThanOp(BuildIntConst(1), BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildLessEqThanOp_ExpectLessEqThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessEqThanOp{}, &IntConst{2}}
	gotVal := BuildLessEqThanOp(BuildIntConst(1), BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildGreaterThanOp_ExpectGreaterThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterThanOp{}, &IntConst{2}}
	gotVal := BuildGreaterThanOp(BuildIntConst(1), BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildGreaterEqThanOp_ExpectGreaterEqThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterEqThanOp{}, &IntConst{2}}
	gotVal := BuildGreaterEqThanOp(BuildIntConst(1), BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsLessThanOp_ExpectLessThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), LESS_THAN_OP, BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsLessEqThanOp_ExpectLessEqThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessEqThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), LESS_EQ_THAN_OP, BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsGreaterThanOp_ExpectGreaterThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), GREATER_THAN_OP, BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsGreaterEqThanOp_ExpectGreaterEqThanOpBetweenIntConstWithOneValueAndIntConstWithTwoValue(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterEqThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), GREATER_EQ_THAN_OP, BuildIntConst(2))

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsNotValid_ExpectNil(t *testing.T) {
	gotVal := BuildCmpOp(BuildIntConst(1), "INVALID_OP", BuildIntConst(2))

	assert.Nil(t, gotVal)
}

func TestBinaryOpKind_ExpectBinaryOpKind(t *testing.T) {
	wantVal := BINARY_OP_CONST
	gotVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestBinaryOpEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	otherVal := BuildIntConst(0)

	assert.False(t, binaryOp.equals(otherVal))
}

func TestBinaryOpEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), neqOperator(), BuildIntConst(1))
	otherVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))

	assert.False(t, binaryOp.equals(otherVal))
}

func TestBinaryOpEquals_ExpectTrue(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	otherVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))

	assert.True(t, binaryOp.equals(otherVal))
}

func TestBuildFunctionCall_ExpectionFunctionCallWithFooNameAndNoArgs(t *testing.T) {
	functionName := "foo"
	functionArgs := []Expr{}

	wantVal := &FunctionCall{&Variable{functionName}, functionArgs}
	gotVal := BuildFunctionCall(BuildVariable(functionName), functionArgs)

	assert.Equal(t, wantVal, gotVal)
}

func TestFunctionCallKind_ExpectFunctionCallKind(t *testing.T) {
	wantVal := FUNCTION_CALL_CONST
	gotVal := BuildFunctionCall(BuildVariable("foo"), []Expr{}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestFunctionCallEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	functionCall := BuildFunctionCall(BuildVariable("foo"), []Expr{})
	otherVal := BuildIntConst(0)

	assert.False(t, functionCall.equals(otherVal))
}

func TestFunctionCallEquals_WhenFalseDueToDiffValues(t *testing.T) {
	functionCall := BuildFunctionCall(BuildVariable("foo"), []Expr{})
	otherVal := BuildFunctionCall(BuildVariable("bar"), []Expr{})

	assert.False(t, functionCall.equals(otherVal))
}

func TestFunctionCallEquals_WhenTrue(t *testing.T) {
	functionCall := BuildFunctionCall(BuildVariable("foo"), []Expr{})
	otherVal := BuildFunctionCall(BuildVariable("foo"), []Expr{})

	assert.True(t, functionCall.equals(otherVal))
}

func TestBuildArray_ExpectArrayConstWithOnlyIntConstWithOneValue(t *testing.T) {
	elems := []Expr{BuildIntConst(1)}

	wantVal := &Array{elems}
	gotVal := BuildArray(elems)

	assert.Equal(t, wantVal, gotVal)
}

func TestArrayKind_ExpectArrayKind(t *testing.T) {
	wantVal := ARRAY_CONST
	gotVal := BuildArray([]Expr{BuildIntConst(1)}).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestArrayEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildIntConst(0)

	assert.False(t, array.equals(otherVal))
}

func TestArrayEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildArray([]Expr{BuildIntConst(1), BuildIntConst(2)})

	assert.False(t, array.equals(otherVal))
}

func TestArrayEquals_ExpectTrue(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildArray([]Expr{BuildIntConst(1)})

	assert.True(t, array.equals(otherVal))
}

func TestBuildTypedExpr_ExpectTypedExprWithIntTypeAndIntConstWithOneValue(t *testing.T) {
	wantVal := &TypedExpr{&IntConst{1}, &IntType{}}
	gotVal := BuildTypedExpr(BuildIntConst(1), BuildIntType())

	assert.Equal(t, wantVal, gotVal)
}

func TestTypedExprKind_ExpectTypedExprKind(t *testing.T) {
	wantVal := TYPED_EXPR
	gotVal := BuildTypedExpr(BuildIntConst(1), BuildIntType()).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestTypedExprEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildIntConst(0)

	assert.False(t, typedExpr.equals(otherVal))
}

func TestTypedExprEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildTypedExpr(BuildIntConst(1), BuildStringType())

	assert.False(t, typedExpr.equals(otherVal))
}

func TestTypedExprEquals_ExpectTrue(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildTypedExpr(BuildIntConst(1), BuildIntType())

	assert.True(t, typedExpr.equals(otherVal))
}

func TestBuildLambda_ExpectLambdaWithTypedExprParamAndBinaryOpBody(t *testing.T) {
	wantVal := &Lambda{
        []Expr{&TypedExpr{&Variable{"foo"}, &StringType{}}}, 
        &BinaryOp{&IntConst{1}, &EqOp{}, &IntConst{1}},
    }
	gotVal := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)

	assert.Equal(t, wantVal, gotVal)
}

func TestLambdaKind_ExpectLambdaKind(t *testing.T) {
	wantVal := LAMBDA_CONST
	gotVal := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	).Kind()

	assert.Equal(t, wantVal, gotVal)
}

func TestLambdaEquals_WhenDiffKinds_ExpectFalse(t *testing.T) {
	lambda := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)
	otherVal := BuildIntConst(0)
	assert.False(t, lambda.equals(otherVal))
}

func TestLambdaEquals_WhenDiffValues_ExpectFalse(t *testing.T) {
	lambda := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)
	otherVal := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("bar"), BuildStringType())},
		BuildBinaryOp(BuildVariable("var#1"), eqOperator(), BuildVariable("var#2")),
	)

	assert.False(t, lambda.equals(otherVal))
}

func TestLambdaEquals_ExpectTrue(t *testing.T) {
	lambda := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)
	otherVal := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)
    
	assert.True(t, lambda.equals(otherVal))
}
