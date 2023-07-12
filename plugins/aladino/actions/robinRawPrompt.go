// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	api_entities "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	converter "github.com/reviewpad/go-lib/converters"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
)

func RobinRawPrompt() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType(), lang.BuildStringType()}, nil),
		Code:           robinRawPromptCode,
		SupportedKinds: []event_processor.TargetEntityKind{event_processor.PullRequest, event_processor.Issue},
	}
}

func robinRawPromptCode(e aladino.Env, args []lang.Value) error {
	target := e.GetTarget()
	targetEntity := target.GetTargetEntity()
	systemPrompt := args[0].(*lang.StringValue).Val
	userPrompt := args[1].(*lang.StringValue).Val
	model := args[2].(*lang.StringValue).Val

	service, ok := e.GetBuiltIns().Services[plugins_aladino_services.ROBIN_SERVICE_KEY]
	if !ok {
		return fmt.Errorf("robin service not found")
	}

	robinClient := service.(api.RobinClient)
	req := &api.RawPromptRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		Token:        e.GetGithubClient().GetToken(),
		Target: &api_entities.TargetEntity{
			Owner:  targetEntity.Owner,
			Repo:   targetEntity.Repo,
			Kind:   converter.ToEntityKind(targetEntity.Kind),
			Number: int32(targetEntity.Number),
		},
		Act:   false,
		Model: model,
	}

	resp, err := robinClient.RawPrompt(e.GetCtx(), req)
	if err != nil {
		return err
	}

	return target.Comment(resp.Reply)
}
