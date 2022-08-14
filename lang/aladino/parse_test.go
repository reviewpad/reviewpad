// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_WhenSingleLine(t *testing.T) {
	input := `$addLabel("small")`
	wantExpr := BuildFunctionCall(
		BuildVariable("addLabel"),
		[]Expr{BuildStringConst("small")},
	)

	gotExpr, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestParse_WhenMultilineLine(t *testing.T) {
	input := `$addLabel("medium multiline")

`
	wantExpr := BuildFunctionCall(
		BuildVariable("addLabel"),
		[]Expr{BuildStringConst("medium multiline")},
	)

	gotExpr, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}

func TestParse_Lambda(t *testing.T) {
	tests := map[string]struct {
		input    string
		wantExpr *Lambda
	}{
		"no arguments": {
			input: `( => 10)`,
			wantExpr: BuildLambda(
				[]Expr{},
				BuildIntConst(10),
			),
		},
		"single argument with wrong type": {
			input: `(1 => 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildIntConst(1)},
				BuildIntConst(10),
			),
		},
		"single argument": {
			input: `($dev => $totalCreatedPullRequests($dev) == 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildVariable("dev")},
				BuildBinaryOp(BuildFunctionCall(BuildVariable("totalCreatedPullRequests"), []Expr{BuildVariable("dev")}), eqOperator(), BuildIntConst(10)),
			),
		},
		"multiple arguments": {
			input: `($a, $b => $a > $b)`,
			wantExpr: BuildLambda(
				[]Expr{BuildVariable("a"), BuildVariable("b")},
				BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
			),
		},
		"nested lambda": {
			input: `($a => ($b => $a > $b))`,
			wantExpr: BuildLambda(
				[]Expr{BuildVariable("a")},
				BuildLambda(
					[]Expr{BuildVariable("b")},
					BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
				),
			),
		},
		"typed expression lambda": {
			input: `($a: Int, $b: Int => $a > $b)`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("a"), BuildIntType()),
					BuildTypedExpr(BuildVariable("b"), BuildIntType())},
				BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotExpr, err := Parse(test.input)
			assert.Nil(t, err)
			assert.Equal(t, test.wantExpr, gotExpr)
		})
	}
}

func TestParse_TypedExpression(t *testing.T) {
	input := `$developer: String`
	wantExpr := BuildTypedExpr(
		BuildVariable("developer"),
		BuildStringType(),
	)

	gotExpr, err := Parse(input)
	assert.Nil(t, err)
	assert.Equal(t, wantExpr, gotExpr)
}
