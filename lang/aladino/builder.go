// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

func BuildFilter(param string, condition Expr) (Expr, error) {
	organizationAST := FunctionCallConstr(
		VariableConstr("organization"),
		[]Expr{},
	)

	ast := FunctionCallConstr(
		VariableConstr("filter"),
		[]Expr{
			organizationAST,
			LambdaConstr(
				[]Expr{TypedExprConstr(VariableConstr(param), BuildStringType())},
				condition,
			),
		},
	)

	return ast, nil
}
