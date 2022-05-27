// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

func buildFilter(param string, condition Expr) (Expr, error) {
	allDevsAST := functionCall(
		variable("allDevs"),
		[]Expr{},
	)

	ast := functionCall(
		variable("filter"),
		[]Expr{
			allDevsAST,
			lambda(
				[]Expr{typedExpr(variable(param), BuildStringType())},
				condition,
			),
		},
	)

	return ast, nil
}
