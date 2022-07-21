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

func TestEqOperator(t *testing.T) {
	wantVal := &EqOp{}
	gotVal := eqOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestNeqOperator(t *testing.T) {
	wantVal := &NeqOp{}
	gotVal := neqOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestAndOperator(t *testing.T) {
	wantVal := &AndOp{}
	gotVal := andOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestOrOperator(t *testing.T) {
	wantVal := &OrOp{}
	gotVal := orOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestLessThanOperator(t *testing.T) {
	wantVal := &LessThanOp{}
	gotVal := lessThanOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestLessEqThanOperator(t *testing.T) {
	wantVal := &LessEqThanOp{}
	gotVal := lessEqThanOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGreaterThanOperator(t *testing.T) {
	wantVal := &GreaterThanOp{}
	gotVal := greaterThanOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGreaterEqThanOperator(t *testing.T) {
	wantVal := &GreaterEqThanOp{}
	gotVal := greaterEqThanOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenEqOp(t *testing.T) {
	wantVal := EQ_OP
	gotVal := eqOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenNeqOp(t *testing.T) {
	wantVal := NEQ_OP
	gotVal := neqOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenAndOp(t *testing.T) {
	wantVal := AND_OP
	gotVal := andOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenOrOp(t *testing.T) {
	wantVal := OR_OP
	gotVal := orOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenLessThanOp(t *testing.T) {
	wantVal := LESS_THAN_OP
	gotVal := lessThanOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenLessEqThanOp(t *testing.T) {
	wantVal := LESS_EQ_THAN_OP
	gotVal := lessEqThanOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenGreaterThanOp(t *testing.T) {
	wantVal := GREATER_THAN_OP
	gotVal := greaterThanOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestGetOperator_WhenGreaterEqThanOp(t *testing.T) {
	wantVal := GREATER_EQ_THAN_OP
	gotVal := greaterEqThanOperator().getOperator()
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildBoolConst(t *testing.T) {
	wantVal := &BoolConst{true}
	gotVal := BuildBoolConst(true)
	assert.Equal(t, wantVal, gotVal)
}

func TestBoolConstKind(t *testing.T) {
	wantVal := BOOL_CONST
	gotVal := BuildBoolConst(true).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestBoolConstEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildIntConst(0)
	assert.False(t, boolConst.equals(otherConst))
}

func TestBoolConstEquals_WhenFalseDueToDiffValues(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildBoolConst(false)
	assert.False(t, boolConst.equals(otherConst))
}

func TestBoolConstEquals_WhenTrue(t *testing.T) {
	boolConst := BuildBoolConst(true)
	otherConst := BuildBoolConst(true)
	assert.True(t, boolConst.equals(otherConst))
}

func TestBuildStringConst(t *testing.T) {
	wantVal := &StringConst{"Lorem Ipsum"}
	gotVal := BuildStringConst("Lorem Ipsum")
	assert.Equal(t, wantVal, gotVal)
}

func TestStringConstKind(t *testing.T) {
	wantVal := STRING_CONST
	gotVal := BuildStringConst("Lorem Ipsum").Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestStringConstEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum")
	otherConst := BuildIntConst(0)
	assert.False(t, stringConst.equals(otherConst))
}

func TestStringConstEquals_WhenFalseDueToDiffValues(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum #1")
	otherConst := BuildStringConst("Lorem Ipsum #2")
	assert.False(t, stringConst.equals(otherConst))
}

func TestStringConstEquals_WhenTrue(t *testing.T) {
	stringConst := BuildStringConst("Lorem Ipsum")
	otherConst := BuildStringConst("Lorem Ipsum")
	assert.True(t, stringConst.equals(otherConst))
}

func TestBuildIntConst(t *testing.T) {
	wantVal := &IntConst{0}
	gotVal := BuildIntConst(0)
	assert.Equal(t, wantVal, gotVal)
}

func TestIntConstKind(t *testing.T) {
	wantVal := INT_CONST
	gotVal := BuildIntConst(0).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestIntConstEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildStringConst("Lorem Ipsum")
	assert.False(t, intConst.equals(otherConst))
}

func TestIntConstEquals_WhenFalseDueToDiffValues(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildIntConst(1)
	assert.False(t, intConst.equals(otherConst))
}

func TestIntConstEquals_WhenTrue(t *testing.T) {
	intConst := BuildIntConst(0)
	otherConst := BuildIntConst(0)
	assert.True(t, intConst.equals(otherConst))
}

func TestBuildRelativeTimeConst_WhenTimeUnitIsYear(t *testing.T) {
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

func TestBuildRelativeTimeConst_WhenTimeUnitIsMonth(t *testing.T) {
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

func TestBuildRelativeTimeConst_WhenTimeUnitIsDay(t *testing.T) {
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

func TestBuildRelativeTimeConst_WhenTimeUnitIsWeek(t *testing.T) {
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

func TestBuildRelativeTimeConst_WhenTimeUnitIsHour(t *testing.T) {
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

func TestBuildRelativeTimeConst_WhenTimeUnitIsMinute(t *testing.T) {
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

func TestBuildTimeConst(t *testing.T) {
	val := "2019-10-12T07:50:52"

	wantVal := &IntConst{
		value: int(time.Date(2019, time.Month(10), 12, 7, 50, 52, 0, time.UTC).Unix()),
	}

	gotVal := BuildTimeConst(val)

	assert.Equal(t, wantVal, gotVal)
}

func TestBuildVariable(t *testing.T) {
	wantVal := &Variable{"test"}
	gotVal := BuildVariable("test")
	assert.Equal(t, wantVal, gotVal)
}

func TestVariableKind(t *testing.T) {
	wantVal := VARIABLE_CONST
	gotVal := BuildVariable("test").Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestVariableEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	variable := BuildVariable("test")
	otherVal := BuildIntConst(0)
	assert.False(t, variable.equals(otherVal))
}

func TestVariableEquals_WhenFalseDueToDiffValues(t *testing.T) {
	variable := BuildVariable("test#1")
	otherVal := BuildVariable("test#2")
	assert.False(t, variable.equals(otherVal))
}

func TestVariableEquals_WhenTrue(t *testing.T) {
	variable := BuildVariable("test")
	otherVal := BuildVariable("test")
	assert.True(t, variable.equals(otherVal))
}

func TestBuildUnaryOp(t *testing.T) {
	wantVal := &UnaryOp{&NotOp{}, &BoolConst{true}}
	gotVal := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildNotOp(t *testing.T) {
	wantVal := &UnaryOp{&NotOp{}, &BoolConst{true}}
	gotVal := BuildNotOp(BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestUnaryOpKind(t *testing.T) {
	wantVal := UNARY_OP_CONST
	gotVal := BuildUnaryOp(notOperator(), BuildBoolConst(true)).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestUnaryOpEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildIntConst(0)
	assert.False(t, unaryOp.equals(otherVal))
}

func TestUnaryOpEquals_WhenFalseDueToDiffValues(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildUnaryOp(notOperator(), BuildBoolConst(false))
	assert.False(t, unaryOp.equals(otherVal))
}

func TestUnaryOpEquals_WhenTrue(t *testing.T) {
	unaryOp := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	otherVal := BuildUnaryOp(notOperator(), BuildBoolConst(true))
	assert.True(t, unaryOp.equals(otherVal))
}

func TestBuildBinaryOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &EqOp{}, &IntConst{1}}
	gotVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildAndOp(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &AndOp{}, &BoolConst{true}}
	gotVal := BuildAndOp(BuildBoolConst(true), BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildOrOp(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &OrOp{}, &BoolConst{true}}
	gotVal := BuildOrOp(BuildBoolConst(true), BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildEqOp(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &EqOp{}, &BoolConst{true}}
	gotVal := BuildEqOp(BuildBoolConst(true), BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildNeqOp(t *testing.T) {
	wantVal := &BinaryOp{&BoolConst{true}, &NeqOp{}, &BoolConst{true}}
	gotVal := BuildNeqOp(BuildBoolConst(true), BuildBoolConst(true))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildLessThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessThanOp{}, &IntConst{2}}
	gotVal := BuildLessThanOp(BuildIntConst(1), BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildLessEqThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessEqThanOp{}, &IntConst{2}}
	gotVal := BuildLessEqThanOp(BuildIntConst(1), BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildGreaterThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterThanOp{}, &IntConst{2}}
	gotVal := BuildGreaterThanOp(BuildIntConst(1), BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildGreaterEqThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterEqThanOp{}, &IntConst{2}}
	gotVal := BuildGreaterEqThanOp(BuildIntConst(1), BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsLessThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), LESS_THAN_OP, BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsLessEqThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &LessEqThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), LESS_EQ_THAN_OP, BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsGreaterThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), GREATER_THAN_OP, BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsGreaterEqThanOp(t *testing.T) {
	wantVal := &BinaryOp{&IntConst{1}, &GreaterEqThanOp{}, &IntConst{2}}
	gotVal := BuildCmpOp(BuildIntConst(1), GREATER_EQ_THAN_OP, BuildIntConst(2))
	assert.Equal(t, wantVal, gotVal)
}

func TestBuildCmpOp_WhenOpIsNotValid(t *testing.T) {
	gotVal := BuildCmpOp(BuildIntConst(1), "INVALID_OP", BuildIntConst(2))
	assert.Nil(t, gotVal)
}

func TestBinaryKind(t *testing.T) {
	wantVal := BINARY_OP_CONST
	gotVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestBinaryOpEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	otherVal := BuildIntConst(0)
	assert.False(t, binaryOp.equals(otherVal))
}

func TestBinaryOpEquals_WhenFalseDueToDiffValues(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), neqOperator(), BuildIntConst(1))
	otherVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	assert.False(t, binaryOp.equals(otherVal))
}

func TestBinaryOpEquals_WhenTrue(t *testing.T) {
	binaryOp := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	otherVal := BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1))
	assert.True(t, binaryOp.equals(otherVal))
}

func TestBuildFunctionCall(t *testing.T) {
	functionName := "foo"
	functionArgs := []Expr{}
	wantVal := &FunctionCall{&Variable{functionName}, functionArgs}
	gotVal := BuildFunctionCall(BuildVariable(functionName), functionArgs)
	assert.Equal(t, wantVal, gotVal)
}

func TestFunctionCallKind(t *testing.T) {
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

func TestBuildArray(t *testing.T) {
	elems := []Expr{BuildIntConst(1)}
	wantVal := &Array{elems}
	gotVal := BuildArray(elems)
	assert.Equal(t, wantVal, gotVal)
}

func TestArrayKind(t *testing.T) {
	wantVal := ARRAY_CONST
	gotVal := BuildArray([]Expr{BuildIntConst(1)}).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestArrayEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildIntConst(0)
	assert.False(t, array.equals(otherVal))
}

func TestArrayEquals_WhenFalseDueToDiffValues(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildArray([]Expr{BuildIntConst(1), BuildIntConst(2)})
	assert.False(t, array.equals(otherVal))
}

func TestArrayEquals_WhenTrue(t *testing.T) {
	array := BuildArray([]Expr{BuildIntConst(1)})
	otherVal := BuildArray([]Expr{BuildIntConst(1)})
	assert.True(t, array.equals(otherVal))
}

func TestBuildTypedExpr(t *testing.T) {
	wantVal := &TypedExpr{&IntConst{1}, &IntType{}}
	gotVal := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	assert.Equal(t, wantVal, gotVal)
}

func TestTypedExprKind(t *testing.T) {
	wantVal := TYPED_EXPR
	gotVal := BuildTypedExpr(BuildIntConst(1), BuildIntType()).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestTypedExprEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildIntConst(0)
	assert.False(t, typedExpr.equals(otherVal))
}

func TestTypedExprEquals_WhenFalseDueToDiffValues(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildTypedExpr(BuildIntConst(1), BuildStringType())
	assert.False(t, typedExpr.equals(otherVal))
}

func TestTypedExprEquals_WhenTrue(t *testing.T) {
	typedExpr := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	otherVal := BuildTypedExpr(BuildIntConst(1), BuildIntType())
	assert.True(t, typedExpr.equals(otherVal))
}

func TestBuildLambda(t *testing.T) {
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

func TestLambdaKind(t *testing.T) {
	wantVal := LAMBDA_CONST
	gotVal := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	).Kind()
	assert.Equal(t, wantVal, gotVal)
}

func TestLambdaEquals_WhenFalseDueToDiffKinda(t *testing.T) {
	lambda := BuildLambda(
		[]Expr{BuildTypedExpr(BuildVariable("foo"), BuildStringType())},
		BuildBinaryOp(BuildIntConst(1), eqOperator(), BuildIntConst(1)),
	)
	otherVal := BuildIntConst(0)
	assert.False(t, lambda.equals(otherVal))
}

func TestLambdaEquals_WhenFalseDueToDiffValues(t *testing.T) {
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

func TestLambdaEquals_WhenTrue(t *testing.T) {
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
