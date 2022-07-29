// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeCheckExec_WhenTypeInferenceFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil)

	expr, err := Parse("$emptyAction(1)")
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

	execExpr, err := TypeCheckExec(mockedEnv, expr)

	assert.Nil(t, execExpr)
	assert.EqualError(t, err, "type inference failed: mismatch in arg types on emptyAction")
}

func TestTypeCheck(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil)

	expr, err := Parse("$emptyAction()")
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

	gotExecExpr, err := TypeCheckExec(mockedEnv, expr)

	wantExecExpr := BuildFunctionCall(BuildVariable("emptyAction"), []Expr{})

	assert.Nil(t, err)
	assert.Equal(t, wantExecExpr, gotExecExpr)
}

func TestTypeCheck_WhenExprIsNotFunctionCall(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil)

	expr, err := Parse("\"not a function call\"")
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

	gotExecExpr, err := TypeCheckExec(mockedEnv, expr)

	assert.Nil(t, gotExecExpr)
	assert.EqualError(t, err, "typecheckexec: StringConst")
}

func TestExec_WhenFunctionArgsEvalFails(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil)

	fc := &FunctionCall{
		name: BuildVariable("invalidCmpOp"),
		arguments: []Expr{
			BuildEqOp(BuildIntConst(0), BuildStringConst("0")),
		},
	}

	err := fc.exec(mockedEnv)

	assert.EqualError(t, err, "eval: left and right operand have different kinds")
}

func TestExec_WhenActionBuiltInNonExisting(t *testing.T) {
	mockedEnv := MockDefaultEnv(t, nil, nil)

	delete(mockedEnv.GetBuiltIns().Actions, "tautology")

	fc := &FunctionCall{
		name: BuildVariable("tautology"),
		arguments: []Expr{
			BuildEqOp(BuildIntConst(1), BuildIntConst(1)),
		},
	}

	err := fc.exec(mockedEnv)

	assert.EqualError(t, err, "exec: tautology not found. are you sure this is a built-in function?")
}

func TestExec_WhenActionIsEnabled(t *testing.T) {
	builtInName := "emptyAction"

	isBuiltInCalled := false

	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			builtInName: {
				Type: BuildFunctionType([]Type{BuildStringType()}, BuildArrayOfType(BuildStringType())),
				Code: func(e Env, args []Value) error {
					isBuiltInCalled = true
					return nil
				},
			},
		},
	}
	mockedEnv := MockDefaultEnvBuiltIns(t, nil, nil, builtIns)

	fc := &FunctionCall{
		name:      BuildVariable(builtInName),
		arguments: []Expr{},
	}

	err := fc.exec(mockedEnv)

	assert.Nil(t, err)
	assert.True(t, isBuiltInCalled)
}

func TestExec_WhenActionIsDisabled(t *testing.T) {
	builtInName := "emptyAction"

	isBuiltInCalled := false

	builtIns := &BuiltIns{
		Actions: map[string]*BuiltInAction{
			builtInName: {
				Type: BuildFunctionType([]Type{BuildStringType()}, BuildArrayOfType(BuildStringType())),
				Code: func(e Env, args []Value) error {
					isBuiltInCalled = true
					return nil
				},
				Disabled: true,
			},
		},
	}
	mockedEnv := MockDefaultEnvBuiltIns(t, nil, nil, builtIns)

	fc := &FunctionCall{
		name:      BuildVariable(builtInName),
		arguments: []Expr{},
	}

	err := fc.exec(mockedEnv)

	assert.Nil(t, err)
	assert.False(t, isBuiltInCalled)
}
