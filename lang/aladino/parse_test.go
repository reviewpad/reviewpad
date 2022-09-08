// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input    string
		wantExpr Expr
	}{
		"single line": {
			input: `$addLabel("small")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("addLabel"),
				[]Expr{BuildStringConst("small")},
			),
		},
		"single line with spaces": {
			input: `$addLabel(  "medium multiline")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("addLabel"),
				[]Expr{BuildStringConst("medium multiline")},
			),
		},
		"typed expression": {
			input: `$developer: String`,
			wantExpr: BuildTypedExpr(
				BuildVariable("developer"),
				BuildStringType(),
			),
		},
		"comment with a string": {
			input: `["hello <a href='https://www.google.com'></a> world", "hello world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello <a href='https://www.google.com'></a> world"),
				BuildStringConst("hello world"),
			}),
		},
		"array of strings": {
			input: `["hello", "world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello"),
				BuildStringConst("world"),
			}),
		},
		"array of typed expression": {
			input: `[$developer: String, "hello", ($dev => $dev == "hello")]`,
			wantExpr: BuildArray([]Expr{
				BuildTypedExpr(
					BuildVariable("developer"),
					BuildStringType(),
				),
				BuildStringConst("hello"),
				BuildLambda(
					[]Expr{BuildVariable("dev")},
					BuildBinaryOp(BuildVariable("dev"), eqOperator(), BuildStringConst("hello")),
				),
			}),
		},
		"lambda no arguments": {
			input: `( => 10)`,
			wantExpr: BuildLambda(
				[]Expr{},
				BuildIntConst(10),
			),
		},
		"lambda single argument with wrong type": {
			input: `(1 => 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildIntConst(1)},
				BuildIntConst(10),
			),
		},
		"lambda single argument": {
			input: `($dev => $totalCreatedPullRequests($dev) == 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildVariable("dev")},
				BuildBinaryOp(BuildFunctionCall(BuildVariable("totalCreatedPullRequests"), []Expr{BuildVariable("dev")}), eqOperator(), BuildIntConst(10)),
			),
		},
		"lambda single argument and operation": {
			input: `($dev => $dev == "hello")`,
			wantExpr: BuildLambda(
				[]Expr{BuildVariable("dev")},
				BuildBinaryOp(BuildVariable("dev"), eqOperator(), BuildStringConst("hello")),
			),
		},
		"lambda multiple arguments": {
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
