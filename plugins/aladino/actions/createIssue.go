// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"github.com/reviewpad/go-lib/entities"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
)

func CreateIssue() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, nil),
		Code:           createIssueCode,
		SupportedKinds: []entities.TargetEntityKind{entities.PullRequest},
	}
}

type InsertIssueInput struct {
	Title          string `graphql:"title" json:"title"`
	ResourceNodeId string `graphql:"resource_node_id" json:"resource_node_id"`
	Results        string `graphql:"results" json:"results"`
	ExecutionId    string `graphql:"execution_id" json:"execution_id"`
}

func (i InsertIssueInput) GetGraphQLType() string {
	return "issues_insert_input"
}

type InsertIssueMutation struct {
	Insert_issue struct {
		Id string
	} `graphql:"insert_issue(object: $object)"`
}

func createIssueCode(e aladino.Env, args []lang.Value) error {
	t := e.GetTarget()
	context := args[0].(*lang.StringValue).Val
	description := args[1].(*lang.StringValue).Val
	log := e.GetLogger().WithField("builtin", "createIssue")

	executionId := e.GetExecutionID()
	nodeId := t.GetTargetEntity().NodeId

	insertInput := InsertIssueInput{
		Title:          context,
		ResourceNodeId: nodeId,
		Results:        description,
		ExecutionId:    executionId,
	}

	err := e.GetNexusClient().WithDebug(true).Mutate(e.GetCtx(), &InsertIssueMutation{}, map[string]interface{}{
		"object": insertInput,
	})

	if err != nil {
		log.WithError(err).Error("Failed to create issue")
		return err
	}

	return nil
}
