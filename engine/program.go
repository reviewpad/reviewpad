// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

type Statement struct {
	code string
}

type Program struct {
	IsFromCommand bool
	statements    []*Statement
}

func BuildStatement(code string) *Statement {
	return &Statement{
		code,
	}
}

func BuildProgram(statements []*Statement, isFromCommand bool) *Program {
	return &Program{
		isFromCommand,
		statements,
	}
}

func (s *Statement) GetStatementCode() string {
	return s.code
}

func (p *Program) GetProgramStatements() []*Statement {
	return p.statements
}

func (program *Program) append(workflowActions []string) {
	for _, workflowAction := range workflowActions {
		statement := BuildStatement(workflowAction)

		program.statements = append(program.statements, statement)
	}
}
