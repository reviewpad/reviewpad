// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/reviewpad/reviewpad/utils/report"
)

type Expr interface {
	Kind() string
	typeinfer(env *TypeEnv) (Type, error)
	Eval(Env) (Value, error)
	equals(Expr) bool
}

const (
	BOOL_CONST          string = "BoolConst"
	INT_CONST           string = "IntConst"
	STRING_CONST        string = "StringConst"
	TIME_CONST          string = "TimeConst"
	VARIABLE_CONST      string = "Variable"
	UNARY_OP_CONST      string = "UnaryOp"
	BINARY_OP_CONST     string = "BinaryOp"
	FUNCTION_CALL_CONST string = "FunctionCall"
	LAMBDA_CONST        string = "Lambda"
	TYPED_EXPR          string = "TypedExpr"
	ARRAY_CONST         string = "Array"
	NOT_OP              string = "!"
	EQ_OP               string = "=="
	NEQ_OP              string = "!="
	AND_OP              string = "&&"
	OR_OP               string = "||"
	LESS_THAN_OP        string = "<"
	LESS_EQ_THAN_OP     string = "<="
	GREATER_THAN_OP     string = ">"
	GREATER_EQ_THAN_OP  string = ">="
)

type UnaryOperator interface {
	getOperator() string
	Eval(exprValue Value) Value
}

type NotOp struct{}

func notOperator() *NotOp { return &NotOp{} }

func (op *NotOp) getOperator() string { return NOT_OP }

type BinaryOperator interface {
	getOperator() string
	Eval(lhs, rhs Value) Value
}

type EqOp struct{}
type NeqOp struct{}
type AndOp struct{}
type OrOp struct{}
type LessThanOp struct{}
type LessEqThanOp struct{}
type GreaterThanOp struct{}
type GreaterEqThanOp struct{}

func eqOperator() *EqOp                       { return &EqOp{} }
func neqOperator() *NeqOp                     { return &NeqOp{} }
func andOperator() *AndOp                     { return &AndOp{} }
func orOperator() *OrOp                       { return &OrOp{} }
func lessThanOperator() *LessThanOp           { return &LessThanOp{} }
func lessEqThanOperator() *LessEqThanOp       { return &LessEqThanOp{} }
func greaterThanOperator() *GreaterThanOp     { return &GreaterThanOp{} }
func greaterEqThanOperator() *GreaterEqThanOp { return &GreaterEqThanOp{} }

func (op *EqOp) getOperator() string            { return EQ_OP }
func (op *NeqOp) getOperator() string           { return NEQ_OP }
func (op *AndOp) getOperator() string           { return AND_OP }
func (op *OrOp) getOperator() string            { return OR_OP }
func (op *LessThanOp) getOperator() string      { return LESS_THAN_OP }
func (op *LessEqThanOp) getOperator() string    { return LESS_EQ_THAN_OP }
func (op *GreaterThanOp) getOperator() string   { return GREATER_THAN_OP }
func (op *GreaterEqThanOp) getOperator() string { return GREATER_EQ_THAN_OP }

type BoolConst struct {
	value bool
}

func boolConst(bval bool) *BoolConst {
	return &BoolConst{value: bval}
}

func trueExpr() *BoolConst {
	return boolConst(true)
}

func falseExpr() *BoolConst {
	return boolConst(false)
}

func (b *BoolConst) Kind() string {
	return BOOL_CONST
}

func (thisBool *BoolConst) equals(other Expr) bool {
	if thisBool.Kind() != other.Kind() {
		return false
	}

	return thisBool.value == other.(*BoolConst).value
}

type StringConst struct {
	value string
}

func stringConst(val string) *StringConst {
	return &StringConst{val}
}

func (c *StringConst) Kind() string {
	return STRING_CONST
}

func (thisString *StringConst) equals(other Expr) bool {
	if thisString.Kind() != other.Kind() {
		return false
	}

	return thisString.value == other.(*StringConst).value
}

type IntConst struct {
	value int
}

func intConst(val int) *IntConst {
	return &IntConst{val}
}

func (i *IntConst) Kind() string {
	return INT_CONST
}

func (thisInt *IntConst) equals(other Expr) bool {
	if thisInt.Kind() != other.Kind() {
		return false
	}

	return thisInt.value == other.(*IntConst).value
}

type TimeConst struct {
	value int64
}

func relativeTimeConst(val string) *TimeConst {
	now := time.Now()

	timeUnitRegex := regexp.MustCompile(`year|month|week|day|hour|minute`)
	timeUnit := timeUnitRegex.FindString(val)

	timeValueRegex := regexp.MustCompile(`^[0-9]+`)
	timeValue, err := strconv.Atoi(timeValueRegex.FindString(val))
	if err != nil {
		log.Fatalf(report.Error(err.Error()))
	}

	switch timeUnit {
	case "year":
		var a = now.AddDate(-timeValue, 0, 0)
		a.UnixMilli()
		return &TimeConst{
			// value:    now.AddDate(-timeValue, 0, 0),
			value: now.AddDate(-timeValue, 0, 0).Unix(),
		}
	case "month":
		return &TimeConst{
			value: now.AddDate(0, -timeValue, 0).Unix(),
		}
	case "day":
		return &TimeConst{
			value: now.AddDate(0, 0, -timeValue).Unix(),
		}
	case "week":
		week := time.Hour * 24 * 7
		return &TimeConst{
			value: now.Add(-week * time.Duration(timeValue)).Unix(),
		}
	case "hour":
		return &TimeConst{
			value: now.Add(-time.Hour * time.Duration(timeValue)).Unix(),
		}
	case "minute":
		return &TimeConst{
			value: now.Add(-time.Minute * time.Duration(timeValue)).Unix(),
		}
	}

	log.Fatalf(report.Error("Unknown time unit %v", timeUnit))
	return &TimeConst{}
}

func timeConst(val string) *TimeConst {
	dateValueRegex := regexp.MustCompile(`^(\d{4})-?(\d{2})-?(\d{2})`)
	dateValue := dateValueRegex.FindSubmatch([]byte(val))

	year, err := strconv.Atoi(string(dateValue[1][:]))
	if err != nil {
		log.Fatalf(report.Error("Error converting year value %q", dateValue[1]))
	}

	month, err := strconv.Atoi(string(dateValue[2][:]))
	if err != nil {
		log.Fatalf(report.Error("Error converting month value %q", dateValue[2]))
	}

	day, err := strconv.Atoi(string(dateValue[3][:]))
	if err != nil {
		log.Fatalf(report.Error("Error converting day value %q", dateValue[3]))
	}

	timeValueRegex := regexp.MustCompile(`T(\d{2}):(\d{2}):(\d{2})$`)
	timeValue := timeValueRegex.FindSubmatch([]byte(val))

	hour := 0
	minute := 0
	second := 0

	if len(timeValue) > 0 {
		hour, err = strconv.Atoi(string(timeValue[1][:]))
		if err != nil {
			log.Fatalf(report.Error("Error converting hour value %q", timeValue[1]))
		}

		minute, err = strconv.Atoi(string(timeValue[2][:]))
		if err != nil {
			log.Fatalf(report.Error("Error converting minute value %q", timeValue[2]))
		}

		second, err = strconv.Atoi(string(timeValue[3][:]))
		if err != nil {
			log.Fatalf(report.Error("Error converting second value %q", timeValue[3]))
		}
	}

	return &TimeConst{
		value: time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC).Unix(),
	}
}

func (i *TimeConst) Kind() string {
	return TIME_CONST
}

func (thisTime *TimeConst) equals(other Expr) bool {
	if thisTime.Kind() != other.Kind() {
		return false
	}

	return thisTime.value == other.(*TimeConst).value
}

type Variable struct {
	ident string
}

func variable(ident string) *Variable {
	return &Variable{ident}
}

func (v *Variable) Kind() string {
	return VARIABLE_CONST
}

func (thisVariable *Variable) equals(other Expr) bool {
	if thisVariable.Kind() != other.Kind() {
		return false
	}

	return thisVariable.ident == other.(*Variable).ident
}

type UnaryOp struct {
	op   UnaryOperator
	expr Expr
}

func unaryOp(op UnaryOperator, expr Expr) *UnaryOp {
	return &UnaryOp{op, expr}
}

func notOp(expr Expr) *UnaryOp { return unaryOp(notOperator(), expr) }

func (b *UnaryOp) Kind() string {
	return UNARY_OP_CONST
}

func (thisUnaryOp *UnaryOp) equals(other Expr) bool {
	if thisUnaryOp.Kind() != other.Kind() {
		return false
	}

	otherUnaryOp := other.(*UnaryOp)
	checkOp := thisUnaryOp.op == otherUnaryOp.op
	exprCheck := thisUnaryOp.expr.equals(otherUnaryOp.expr)

	return checkOp && exprCheck
}

type BinaryOp struct {
	lhs Expr
	op  BinaryOperator
	rhs Expr
}

func binaryOp(lhs Expr, op BinaryOperator, rhs Expr) *BinaryOp {
	return &BinaryOp{lhs, op, rhs}
}

func andOp(lhs Expr, rhs Expr) *BinaryOp         { return binaryOp(lhs, andOperator(), rhs) }
func orOp(lhs Expr, rhs Expr) *BinaryOp          { return binaryOp(lhs, orOperator(), rhs) }
func eqOp(lhs Expr, rhs Expr) *BinaryOp          { return binaryOp(lhs, eqOperator(), rhs) }
func neqOp(lhs Expr, rhs Expr) *BinaryOp         { return binaryOp(lhs, neqOperator(), rhs) }
func lessThanOp(lhs Expr, rhs Expr) *BinaryOp    { return binaryOp(lhs, lessThanOperator(), rhs) }
func lessEqThanOp(lhs Expr, rhs Expr) *BinaryOp  { return binaryOp(lhs, lessEqThanOperator(), rhs) }
func greaterThanOp(lhs Expr, rhs Expr) *BinaryOp { return binaryOp(lhs, greaterThanOperator(), rhs) }
func greaterEqThanOp(lhs Expr, rhs Expr) *BinaryOp {
	return binaryOp(lhs, greaterEqThanOperator(), rhs)
}

func cmpOp(lhs Expr, op string, rhs Expr) Expr {
	switch op {
	case LESS_THAN_OP:
		return lessThanOp(lhs, rhs)
	case LESS_EQ_THAN_OP:
		return lessEqThanOp(lhs, rhs)
	case GREATER_THAN_OP:
		return greaterThanOp(lhs, rhs)
	case GREATER_EQ_THAN_OP:
		return greaterEqThanOp(lhs, rhs)
	default:
		fmt.Printf("cmpOp: invalid op %v\n", op)
		return nil
	}
}

func (b *BinaryOp) Kind() string {
	return BINARY_OP_CONST
}

func (thisBinOp *BinaryOp) equals(other Expr) bool {
	if thisBinOp.Kind() != other.Kind() {
		return false
	}

	otherBinaryOp := other.(*BinaryOp)
	checkOp := thisBinOp.op == otherBinaryOp.op
	lhsCheck := thisBinOp.lhs.equals(otherBinaryOp.lhs)
	rhsCheck := thisBinOp.rhs.equals(otherBinaryOp.rhs)

	return checkOp && lhsCheck && rhsCheck
}

type FunctionCall struct {
	name      *Variable
	arguments []Expr
}

func functionCall(name *Variable, arguments []Expr) *FunctionCall {
	return &FunctionCall{name, arguments}
}

func (fc *FunctionCall) Kind() string {
	return FUNCTION_CALL_CONST
}

func (thisFnCall *FunctionCall) equals(other Expr) bool {
	if thisFnCall.Kind() != other.Kind() {
		return false
	}

	otherFunctionCall := other.(*FunctionCall)
	checkFunctionName := thisFnCall.name.equals(otherFunctionCall.name)
	checkArgs := EqualList(thisFnCall.arguments, otherFunctionCall.arguments)

	return checkFunctionName && checkArgs
}

type Array struct {
	elems []Expr
}

func array(elems []Expr) *Array {
	return &Array{elems}
}

func (a *Array) Kind() string {
	return ARRAY_CONST
}

func (thisArray *Array) equals(other Expr) bool {
	if thisArray.Kind() != other.Kind() {
		return false
	}

	otherArray := other.(*Array)

	return EqualList(thisArray.elems, otherArray.elems)
}

func EqualList(left []Expr, right []Expr) bool {
	if len(left) != len(right) {
		return false
	}

	for i, lExpr := range left {
		rExpr := right[i]
		if !lExpr.equals(rExpr) {
			return false
		}
	}

	return true
}

type TypedExpr struct {
	expr   Expr
	typeOf Type
}

func typedExpr(expr Expr, typeOf Type) *TypedExpr {
	return &TypedExpr{expr, typeOf}
}

func (te *TypedExpr) Kind() string {
	return TYPED_EXPR
}

func (te *TypedExpr) equals(other Expr) bool {
	if te.Kind() != other.Kind() {
		return false
	}

	otherTypedExpr := other.(*TypedExpr)

	return te.expr.equals(otherTypedExpr.expr) &&
		te.typeOf.equals(otherTypedExpr.typeOf)
}

type Lambda struct {
	parameters []Expr
	body       Expr
}

func lambda(parameters []Expr, body Expr) *Lambda {
	return &Lambda{parameters, body}
}

func (l *Lambda) Kind() string {
	return LAMBDA_CONST
}

func (thisLambda *Lambda) equals(other Expr) bool {
	if thisLambda.Kind() != other.Kind() {
		return false
	}

	otherLambda := other.(*Lambda)
	checkBody := thisLambda.body.equals(otherLambda.body)
	checkParameters := EqualList(thisLambda.parameters, otherLambda.parameters)

	return checkBody && checkParameters
}
