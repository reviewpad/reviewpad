// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestBuildFilter(t *testing.T) {
	param := "test"
	condition, _ := aladino.Parse("$totalCreatedPRs($dev) < 10")

	wantAST := aladino.FunctionCallConstr(
		aladino.VariableConstr("filter"),
		[]aladino.Expr{
			aladino.FunctionCallConstr(
				aladino.VariableConstr("organization"),
				[]aladino.Expr{},
			),
			aladino.LambdaConstr(
				[]aladino.Expr{
					aladino.TypedExprConstr(
						aladino.VariableConstr(param),
						aladino.BuildStringType()),
				},
				condition,
			),
		},
	)

	gotAST, err := aladino.BuildFilter(param, condition)

	assert.Nil(t, err)
	assert.Equal(t, wantAST, gotAST)
}
