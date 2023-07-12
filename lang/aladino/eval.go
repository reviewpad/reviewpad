// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
	"strings"

	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
)

type UnsupportedKindError struct {
	Kind           event_processor.TargetEntityKind
	SupportedKinds []event_processor.TargetEntityKind
	BuiltIn        string
}

func (e *UnsupportedKindError) Error() string {
	return fmt.Sprintf("unsupported kind %v for built-in %v", e.Kind, e.BuiltIn)
}

func (e *UnsupportedKindError) SupportedOn() string {
	supportedOn := []string{}
	for _, kind := range e.SupportedKinds {
		supportedOn = append(supportedOn, kind.String())
	}

	return strings.Join(supportedOn, ", ")
}

func (u *UnaryOp) Eval(e Env) (lang.Value, error) {
	exprValue, exprErr := u.expr.Eval(e)
	if exprErr != nil {
		return nil, exprErr
	}

	operator := u.op

	return operator.Eval(exprValue), nil
}

func (b *BinaryOp) Eval(e Env) (lang.Value, error) {
	leftValue, leftErr := b.lhs.Eval(e)
	if leftErr != nil {
		return nil, leftErr
	}

	rightValue, rightErr := b.rhs.Eval(e)
	if rightErr != nil {
		return nil, rightErr
	}

	if !leftValue.HasKindOf(rightValue.Kind()) {
		return nil, fmt.Errorf("eval: left and right operand have different kinds")
	}

	operator := b.op

	return operator.Eval(leftValue, rightValue), nil
}

func (v *Variable) Eval(e Env) (lang.Value, error) {
	variableName := v.ident

	if val, ok := e.GetRegisterMap()[variableName]; ok {
		return val, nil
	}

	// TODO: we need to have a lint step to check if a variable name is a built-in function name
	variable, ok := e.GetRegisterMap()[fmt.Sprintf("@variable:$%s", variableName)]
	if ok {
		return variable, nil
	}

	fn, ok := e.GetBuiltIns().Functions[variableName]
	if !ok {
		return nil, fmt.Errorf("eval: failure on %v", variableName)
	}

	entityKind := e.GetTarget().GetTargetEntity().Kind

	for _, supportedKind := range fn.SupportedKinds {
		if entityKind == supportedKind {
			return fn.Code(e, []lang.Value{})
		}
	}

	return nil, fmt.Errorf("eval: unsupported kind %v", entityKind)
}

func (b *BoolConst) Eval(e Env) (lang.Value, error) {
	return lang.BuildBoolValue(b.value), nil
}

func (c *StringConst) Eval(e Env) (lang.Value, error) {
	return lang.BuildStringValue(c.value), nil
}

func (i *IntConst) Eval(e Env) (lang.Value, error) {
	return lang.BuildIntValue(i.value), nil
}

func (fc *FunctionCall) Eval(e Env) (lang.Value, error) {
	args := make([]lang.Value, len(fc.arguments))
	for i, elem := range fc.arguments {
		value, err := elem.Eval(e)

		if err != nil {
			return nil, err
		}

		args[i] = value
	}

	fn, ok := e.GetBuiltIns().Functions[fc.name.ident]
	if !ok {
		return nil, fmt.Errorf("eval: failure on %v", fc.name.ident)
	}

	entityKind := e.GetTarget().GetTargetEntity().Kind

	for _, supportedKind := range fn.SupportedKinds {
		if entityKind == supportedKind {
			collectedData := map[string]interface{}{
				"builtin": fc.name.ident,
			}

			if err := e.GetCollector().Collect("Ran Builtin", collectedData); err != nil {
				e.GetLogger().Errorf("error collection built-in run: %v\n", err)
			}

			return fn.Code(e, args)
		}
	}

	return nil, &UnsupportedKindError{
		Kind:           entityKind,
		SupportedKinds: fn.SupportedKinds,
		BuiltIn:        fc.name.ident,
	}
}

func (lambda *Lambda) Eval(e Env) (lang.Value, error) {
	fn := func(args []lang.Value) lang.Value {
		// TODO: We need to deep copy the register map and use in the eval at L79
		for i, elem := range lambda.parameters {
			paramIdent := elem.(*TypedExpr).expr.(*Variable).ident

			e.GetRegisterMap()[paramIdent] = args[i]
		}

		fnVal, err := lambda.body.Eval(e)
		if err != nil {
			return nil
		}

		return fnVal
	}

	return lang.BuildFunctionValue(fn), nil
}

func (te *TypedExpr) Eval(e Env) (lang.Value, error) {
	return te.expr.Eval(e)
}

func (a *Array) Eval(e Env) (lang.Value, error) {
	values := make([]lang.Value, len(a.elems))
	for i, elem := range a.elems {
		value, err := elem.Eval(e)

		if err != nil {
			return nil, err
		}

		values[i] = value
	}

	return lang.BuildArrayValue(values), nil
}

func Eval(env Env, expr Expr) (lang.Value, error) {
	val, err := expr.Eval(env)

	if err != nil {
		return nil, err
	}

	return val, nil
}

// EvalCondition evaluates a boolean expression
// Pre-condition: the type of expr is BoolType
func EvalCondition(env Env, expr Expr) (bool, error) {
	boolVal, err := expr.Eval(env)

	if err != nil {
		return false, err
	}

	return boolVal.(*lang.BoolValue).Val, nil
}

func (op *NotOp) Eval(exprVal lang.Value) lang.Value {
	return lang.BuildBoolValue(!exprVal.(*lang.BoolValue).Val)
}

func (op *EqOp) Eval(lhs, rhs lang.Value) lang.Value {
	return lang.BuildBoolValue(lhs.Equals(rhs))
}

func (op *NeqOp) Eval(lhs, rhs lang.Value) lang.Value {
	return lang.BuildBoolValue(!lhs.Equals(rhs))
}

func (op *AndOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.BoolValue).Val
	rightValue := rhs.(*lang.BoolValue).Val

	return lang.BuildBoolValue(leftValue && rightValue)
}

func (op *OrOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.BoolValue).Val
	rightValue := rhs.(*lang.BoolValue).Val

	return lang.BuildBoolValue(leftValue || rightValue)
}

func (op *LessThanOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.IntValue).Val
	rightValue := rhs.(*lang.IntValue).Val

	return lang.BuildBoolValue(leftValue < rightValue)
}

func (op *LessEqThanOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.IntValue).Val
	rightValue := rhs.(*lang.IntValue).Val

	return lang.BuildBoolValue(leftValue <= rightValue)
}

func (op *GreaterThanOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.IntValue).Val
	rightValue := rhs.(*lang.IntValue).Val

	return lang.BuildBoolValue(leftValue > rightValue)
}

func (op *GreaterEqThanOp) Eval(lhs, rhs lang.Value) lang.Value {
	leftValue := lhs.(*lang.IntValue).Val
	rightValue := rhs.(*lang.IntValue).Val

	return lang.BuildBoolValue(leftValue >= rightValue)
}
