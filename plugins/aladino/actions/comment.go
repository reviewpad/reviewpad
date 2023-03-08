// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	pbe "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
)

func Comment() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, nil),
		Code:           commentCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func commentCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)
	pullRequest := t.PullRequest
	repo := pullRequest.GetBase().GetRepo()
	commentBody := args[0].(*aladino.StringValue).Val

	service, ok := e.GetBuiltIns().Services[plugins_aladino_services.CODEHOST_SERVICE_KEY]
	if !ok {
		return fmt.Errorf("code host service not found")
	}

	codehostClient := service.(api.HostsClient)
	req := &api.PostGeneralCommentRequest{
		Host:             pbe.Host_GITHUB,
		HostUri:          "https://github.com",
		Slug:             repo.GetFullName(),
		ExternalRepoId:   fmt.Sprint(repo.GetID()),
		ExternalReviewId: fmt.Sprint(pullRequest.GetID()),
		ReviewNumber:     int32(pullRequest.GetNumber()),
		AccessToken:      e.GetGithubClient().GetToken(),
		Comment: &pbe.ReviewComment{
			Body: commentBody,
		},
	}

	reply, err := codehostClient.PostGeneralComment(e.GetCtx(), req)
	if err != nil {
		return err
	}

	e.GetLogger().Infof("%v\n", reply.Comment)
	return nil
}
