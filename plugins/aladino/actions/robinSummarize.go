// Copyright 2023 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	api_entities "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	converter "github.com/reviewpad/go-lib/converters"
	"github.com/reviewpad/go-lib/entities"
	lib_http "github.com/reviewpad/go-lib/http"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
)

func RobinSummarize() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:              lang.BuildFunctionType([]lang.Type{lang.BuildStringType(), lang.BuildStringType()}, nil),
		Code:              robinSummarizeCode,
		SupportedKinds:    []entities.TargetEntityKind{entities.PullRequest, entities.Issue},
		RunAsynchronously: true,
	}
}

func robinSummarizeCode(e aladino.Env, args []lang.Value) error {
	target := e.GetTarget()
	targetEntity := target.GetTargetEntity()
	summaryMode := args[0].(*lang.StringValue).Val
	model := args[1].(*lang.StringValue).Val

	log := e.GetLogger().WithField("builtin", "summarize")

	service, ok := e.GetBuiltIns().Services[plugins_aladino_services.ROBIN_SERVICE_KEY]
	if !ok {
		return fmt.Errorf("robin service not found")
	}

	robinClient := service.(api.RobinClient)
	req := &api.SummarizeRequest{
		Token: e.GetGithubClient().GetToken(),
		Target: &api_entities.TargetEntity{
			Owner:  targetEntity.Owner,
			Repo:   targetEntity.Repo,
			Kind:   converter.ToEntityKind(targetEntity.Kind),
			Number: int32(targetEntity.Number),
		},
		Act:      false,
		Extended: summaryMode == "extended",
		Model:    model,
	}

	resp, err := robinClient.Summarize(e.GetCtx(), req)
	if err != nil {
		log.Infof("summarize failed with %v", err)
		return fmt.Errorf("summarize failed on request %v - please contact us on [Discord](https://reviewpad.com/discord)", lib_http.InRequestID(e.GetCtx()))
	}

	return target.Comment(fmt.Sprintf("**AI-Generated Summary**: %s", resp.Summary))
}
