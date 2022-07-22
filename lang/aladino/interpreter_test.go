// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine"
	"github.com/stretchr/testify/assert"
)

func TestBuildGroupAST_WhenGroupTypeFilterIsSetAndParseFails(t *testing.T) {
	groupName := "senior-developers"

	gotExpr, err := buildGroupAST(
		engine.GroupTypeFilter,
		"$group(\""+groupName+"\")",
		"dev",
		"$hasFileExtensions(",
	)

	assert.Nil(t, gotExpr)
	assert.EqualError(t, err, "buildGroupAST: parse error: failed to build AST on input $hasFileExtensions(")
}

func TestBuildGroupAST_WhenGroupTypeFilterIsSet(t *testing.T) {
	groupName := "senior-developers"

	gotExpr, err := buildGroupAST(
		engine.GroupTypeFilter,
		"$group(\""+groupName+"\")",
		"dev",
		"$hasFileExtensions([\".ts\"])",
	)

	wantExpr := BuildFunctionCall(
		BuildVariable("filter"),
		[]Expr{
			BuildFunctionCall(
				BuildVariable("organization"),
				[]Expr{},
			),
			BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildFunctionCall(
					BuildVariable("hasFileExtensions"),
					[]Expr{
						BuildStringConst(".ts"),
					},
				),
			),
		},
	)

	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSet(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"[\"jane\"]",
		"",
		"",
	)

	gotVal := mockedEnv.GetRegisterMap()[groupName]

	wantVal := BuildArrayValue([]Value{
		BuildStringValue("jane"),
	})

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSetAndEvalGroupTypeInferenceFails(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"$group(1)",
		"",
		"",
	)

	assert.EqualError(t, err, "type inference failed: mismatch in arg types on group")
}

func TestProcessGroup_WhenGroupTypeFilterIsNotSetAndGroupExprIsNotAnArray(t *testing.T) {
	mockedEnv, err := MockDefaultEnvWithBuiltIns(nil, nil, plugins_aladino.PluginBuiltIns())
	if err != nil {
		log.Fatalf("MockDefaultEnvWithBuiltIns failed: %v", err)
	}

	mockedInterpreter := &Interpreter{
		Env: mockedEnv,
	}

	groupName := "senior-developers"

	err = mockedInterpreter.ProcessGroup(
		groupName,
		engine.GroupKindDeveloper,
		engine.GroupTypeStatic,
		"1",
		"",
		"",
	)

	assert.EqualError(t, err, "expression is not a valid group")
}
