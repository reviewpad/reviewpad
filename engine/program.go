// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

type Metadata struct {
	workflow    PadWorkflow
	triggeredBy []PadWorkflowRule
}

type Statement struct {
	code     string
	metadata *Metadata
}

type Program struct {
	statements []*Statement
}

func BuildMetadata(workflow PadWorkflow, triggeredBy []PadWorkflowRule) *Metadata {
    return &Metadata{
        workflow,
        triggeredBy,
    }
}

func BuildStatement(code string, metadata *Metadata) *Statement {
    return &Statement{
        code,
        metadata,
    }
}

func BuildProgram(statements []*Statement) *Program {
    return &Program{
        statements,
    }
}

func (m *Metadata) GetMetadataWorkflow() PadWorkflow {
    return m.workflow
}

func (m *Metadata) GetMetadataTriggeredBy() []PadWorkflowRule {
    return m.triggeredBy
}

func (s *Statement) GetStatementCode() string {
    return s.code
}

func (s *Statement) GetStatementMetadata() *Metadata {
    return s.metadata
}

func (p *Program) GetProgramStatements() []*Statement {
    return p.statements
}

func (program *Program) append(workflowActions []string, workflow PadWorkflow, workflowRules []PadWorkflowRule) {
	for _, workflowAction := range workflowActions {
        metadata := BuildMetadata(workflow, workflowRules)
        statement := BuildStatement(workflowAction, metadata)

		program.statements = append(program.statements, statement)
	}
}
