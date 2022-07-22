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
	mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    expr, err := Parse("$addLabel(1)")
    if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

    execExpr, err := TypeCheckExec(mockedEnv, expr)

    assert.Nil(t, execExpr)
    assert.EqualError(t, err, "type inference failed: mismatch in arg types on addLabel")
}

func TestTypeCheck(t *testing.T) {
    mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    expr, err := Parse("$addLabel(\"label\")")
    if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

    gotExecExpr, err := TypeCheckExec(mockedEnv, expr)

    wantExecExpr := BuildFunctionCall(BuildVariable("addLabel"), []Expr{BuildStringConst("label")})

    assert.Nil(t, err)
    assert.Equal(t, wantExecExpr, gotExecExpr)
}

func TestTypeCheck_WhenExprIsNotFunctionCall(t *testing.T) {
    mockedEnv, err := MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

    expr, err := Parse("\"not a function call\"")
    if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

    gotExecExpr, err := TypeCheckExec(mockedEnv, expr)

    assert.Nil(t, gotExecExpr)
    assert.EqualError(t, err, "typecheckexec: StringConst")
}
