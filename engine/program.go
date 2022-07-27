// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

type Metadata struct {
	Workflow    PadWorkflow
	TriggeredBy []PadWorkflowRule
}

type Statement struct {
	Code     string
	Metadata *Metadata
}

type Program struct {
	Statements []*Statement
}

func (program *Program) append(workflowActions []string, workflow PadWorkflow, workflowRules []PadWorkflowRule) {
	for _, action := range workflowActions {
		statement := &Statement{
			Code: action,
			Metadata: &Metadata{
				Workflow:    workflow,
				TriggeredBy: workflowRules,
			},
		}

		program.Statements = append(program.Statements, statement)
	}
}
