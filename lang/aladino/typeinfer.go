// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v4/lang"
)

func TypeInference(e Env, expr Expr) (lang.Type, error) {
	return expr.typeinfer(NewTypeEnv(e))
}

func typesinfer(env TypeEnv, exprs []Expr) ([]lang.Type, error) {
	exprsTy := make([]lang.Type, len(exprs))
	for i, expr := range exprs {
		exprTy, err := expr.typeinfer(env)
		if err != nil {
			return nil, err
		}

		exprsTy[i] = exprTy
	}

	return exprsTy, nil
}

func (u *UnaryOp) typeinfer(env TypeEnv) (lang.Type, error) {
	exprType, exprErr := u.expr.typeinfer(env)
	if exprErr != nil {
		return nil, exprErr
	}

	switch u.op.getOperator() {
	case NOT_OP:
		if exprType.Kind() == lang.BOOL_TYPE {
			return lang.BuildBoolType(), nil
		}
	}
	return nil, fmt.Errorf("type inference failed")
}

func (b *BinaryOp) typeinfer(env TypeEnv) (lang.Type, error) {
	lhsType, errLeft := b.lhs.typeinfer(env)
	if errLeft != nil {
		return nil, errLeft
	}

	rhsType, errRight := b.rhs.typeinfer(env)
	if errRight != nil {
		return nil, errRight
	}

	switch b.op.getOperator() {
	case EQ_OP, NEQ_OP:
		if lhsType.Equals(rhsType) {
			return lang.BuildBoolType(), nil
		}
	case GREATER_EQ_THAN_OP, GREATER_THAN_OP, LESS_EQ_THAN_OP, LESS_THAN_OP:
		if lhsType.Equals(lang.BuildIntType()) && rhsType.Equals(lang.BuildIntType()) {
			return lang.BuildBoolType(), nil
		}
	case AND_OP, OR_OP:
		if lhsType.Equals(lang.BuildBoolType()) && rhsType.Equals(lang.BuildBoolType()) {
			return lang.BuildBoolType(), nil
		}
	}

	return nil, fmt.Errorf("type inference failed")
}

func (fc *FunctionCall) typeinfer(env TypeEnv) (lang.Type, error) {
	argsTy, err := typesinfer(env, fc.arguments)
	if err != nil {
		return nil, err
	}

	fcType, err := fc.name.typeinfer(env)
	if err != nil {
		return nil, err
	}

	ty := fcType.(*lang.FunctionType)
	if lang.Equals(argsTy, ty.ParamTypes()) {
		return ty.ReturnType(), nil
	}

	return nil, fmt.Errorf("type inference failed: mismatch in arg types on %v", fc.name.ident)
}

func (l *Lambda) typeinfer(env TypeEnv) (lang.Type, error) {
	paramsTy, err := typesinfer(env, l.parameters)
	if err != nil {
		return nil, err
	}

	bodyType, err := l.body.typeinfer(env)
	if err != nil {
		return nil, err
	}

	return lang.BuildFunctionType(paramsTy, bodyType), nil
}

func (te *TypedExpr) typeinfer(env TypeEnv) (lang.Type, error) {
	if te.expr.Kind() != VARIABLE_CONST {
		return nil, fmt.Errorf("typed expression %v is not a variable", te.expr)
	}

	varIdent := te.expr.(*Variable).ident
	env[varIdent] = te.typeOf

	return te.typeOf, nil
}

// TODO: Fix variable shadowing
func (v *Variable) typeinfer(env TypeEnv) (lang.Type, error) {
	varName := v.ident
	varType, ok := env[varName]
	if !ok {
		// if the built-in is not found in the environment
		// we check if it's a variable
		varType, ok := env[fmt.Sprintf("@variable:$%s", varName)]
		if ok {
			return varType, nil
		}

		return nil, fmt.Errorf("no type for built-in %v", varName)
	}

	return varType, nil
}

func (c *StringConst) typeinfer(env TypeEnv) (lang.Type, error) {
	return lang.BuildStringType(), nil
}

func (i *IntConst) typeinfer(env TypeEnv) (lang.Type, error) {
	return lang.BuildIntType(), nil
}

func (b *BoolConst) typeinfer(env TypeEnv) (lang.Type, error) {
	return lang.BuildBoolType(), nil
}

func (a *Array) typeinfer(env TypeEnv) (lang.Type, error) {
	elemsTy, err := typesinfer(env, a.elems)
	if err != nil {
		return nil, err
	}

	return lang.BuildArrayType(elemsTy), nil
}
