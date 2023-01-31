// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"
)

func (u *UnaryOp) Eval(e Env) (Value, error) {
	exprValue, exprErr := u.expr.Eval(e)
	if exprErr != nil {
		return nil, exprErr
	}

	operator := u.op

	return operator.Eval(exprValue), nil
}

func (b *BinaryOp) Eval(e Env) (Value, error) {
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

func (v *Variable) Eval(e Env) (Value, error) {
	variableName := v.ident

	if val, ok := e.GetRegisterMap()[variableName]; ok {
		return val, nil
	}

	fn, ok := e.GetBuiltIns().Functions[variableName]
	if !ok {
		return nil, fmt.Errorf("eval: failure on %v", variableName)
	}

	entityKind := e.GetTarget().GetTargetEntity().Kind

	for _, supportedKind := range fn.SupportedKinds {
		if entityKind == supportedKind {
			collectedData := map[string]interface{}{
				"builtin": variableName,
			}

			if err := e.GetCollector().Collect("Ran Builtin", collectedData); err != nil {
				e.GetLogger().Errorf("error collection built-in run: %v\n", err)
			}

			return fn.Code(e, []Value{})
		}
	}

	return nil, fmt.Errorf("eval: unsupported kind %v", entityKind)
}

func (b *BoolConst) Eval(e Env) (Value, error) {
	return BuildBoolValue(b.value), nil
}

func (c *StringConst) Eval(e Env) (Value, error) {
	return BuildStringValue(c.value), nil
}

func (i *IntConst) Eval(e Env) (Value, error) {
	return BuildIntValue(i.value), nil
}

func (fc *FunctionCall) Eval(e Env) (Value, error) {
	args := make([]Value, len(fc.arguments))
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

	return nil, fmt.Errorf("eval: unsupported kind %v", entityKind)
}

func (lambda *Lambda) Eval(e Env) (Value, error) {
	fn := func(args []Value) Value {
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

	return BuildFunctionValue(fn), nil
}

func (te *TypedExpr) Eval(e Env) (Value, error) {
	return te.expr.Eval(e)
}

func (a *Array) Eval(e Env) (Value, error) {
	values := make([]Value, len(a.elems))
	for i, elem := range a.elems {
		value, err := elem.Eval(e)

		if err != nil {
			return nil, err
		}

		values[i] = value
	}

	return BuildArrayValue(values), nil
}

func Eval(env Env, expr Expr) (Value, error) {
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

	return boolVal.(*BoolValue).Val, nil
}

func (op *NotOp) Eval(exprVal Value) Value {
	return BuildBoolValue(!exprVal.(*BoolValue).Val)
}

func (op *EqOp) Eval(lhs, rhs Value) Value {
	return BuildBoolValue(lhs.Equals(rhs))
}

func (op *NeqOp) Eval(lhs, rhs Value) Value {
	return BuildBoolValue(!lhs.Equals(rhs))
}

func (op *AndOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*BoolValue).Val
	rightValue := rhs.(*BoolValue).Val

	return BuildBoolValue(leftValue && rightValue)
}

func (op *OrOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*BoolValue).Val
	rightValue := rhs.(*BoolValue).Val

	return BuildBoolValue(leftValue || rightValue)
}

func (op *LessThanOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*IntValue).Val
	rightValue := rhs.(*IntValue).Val

	return BuildBoolValue(leftValue < rightValue)
}

func (op *LessEqThanOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*IntValue).Val
	rightValue := rhs.(*IntValue).Val

	return BuildBoolValue(leftValue <= rightValue)
}

func (op *GreaterThanOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*IntValue).Val
	rightValue := rhs.(*IntValue).Val

	return BuildBoolValue(leftValue > rightValue)
}

func (op *GreaterEqThanOp) Eval(lhs, rhs Value) Value {
	leftValue := lhs.(*IntValue).Val
	rightValue := rhs.(*IntValue).Val

	return BuildBoolValue(leftValue >= rightValue)
}
