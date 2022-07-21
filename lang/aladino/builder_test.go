// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestBuildFilter(t *testing.T) {
	param := "test"
	condition, _ := aladino.Parse("$totalCreatedPRs($dev) < 10")

	wantAST := aladino.BuildFunctionCall(
		aladino.BuildVariable("filter"),
		[]aladino.Expr{
			aladino.BuildFunctionCall(
				aladino.BuildVariable("organization"),
				[]aladino.Expr{},
			),
			aladino.BuildLambda(
				[]aladino.Expr{
					aladino.BuildTypedExpr(
						aladino.BuildVariable(param),
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
