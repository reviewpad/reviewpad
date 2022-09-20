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
		"comment with a string": {
			input: `["hello <a href='https://www.google.com'></a> world", "hello world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello <a href='https://www.google.com'></a> world"),
				BuildStringConst("hello world"),
			}),
		},
		"comment a string with spaced strings": {
			input: `$comment("hello \"world\" and \"world\" again")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("comment"),
				[]Expr{BuildStringConst("hello \"world\" and \"world\" again")},
			),
		},
		"array of strings": {
			input: `["hello", "world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello"),
				BuildStringConst("world"),
			}),
		},
		"array of typed expression": {
			input: `[($developer: String => $developer), "hello", ($dev: String => $dev == "hello")]`,
			wantExpr: BuildArray([]Expr{
				BuildLambda(
					[]Expr{
						BuildTypedExpr(
							BuildVariable("developer"),
							BuildStringType(),
						),
					},
					BuildVariable("developer"),
				),
				BuildStringConst("hello"),
				BuildLambda(
					[]Expr{
						BuildTypedExpr(
							BuildVariable("dev"),
							BuildStringType(),
						),
					},
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
		"lambda single argument": {
			input: `($dev: String => $totalCreatedPullRequests($dev) == 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildBinaryOp(BuildFunctionCall(BuildVariable("totalCreatedPullRequests"), []Expr{BuildVariable("dev")}), eqOperator(), BuildIntConst(10)),
			),
		},
		"lambda single argument and operation": {
			input: `($dev: String => $dev == "hello")`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildBinaryOp(BuildVariable("dev"), eqOperator(), BuildStringConst("hello")),
			),
		},
		"lambda multiple arguments": {
			input: `($a: Int, $b: Int => $a > $b)`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildIntType()),
					BuildTypedExpr(BuildVariable("b"), BuildIntType()),
				},
				BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
			),
		},
		"higher order functions": {
			input: `$any($reviewers(), ($dev: String => $isElementOf($dev, $team("security"))))`,
			wantExpr: BuildFunctionCall(
				BuildVariable("any"),
				[]Expr{
					BuildFunctionCall(
						BuildVariable("reviewers"),
						[]Expr{},
					),
					BuildLambda(
						[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildVariable("dev"),
								BuildFunctionCall(
									BuildVariable("team"),
									[]Expr{BuildStringConst("security")},
								),
							},
						),
					),
				},
			),
		},
		"nested lambda": {
			input: `($a: Int => ($b: Int => $a > $b))`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("a"), BuildIntType())},
				BuildLambda(
					[]Expr{BuildTypedExpr(BuildVariable("b"), BuildIntType())},
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
