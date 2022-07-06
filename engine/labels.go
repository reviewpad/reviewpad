// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/utils"
)

func validateLabelColor(label *PadLabel) error {
	if label.Color != "" {
		matched, _ := regexp.MatchString(`(?i)^([0-9A-F]{6}){1,2}$`, label.Color)
		if !matched {
			if strings.HasPrefix(label.Color, "#") {
				return execError("evalLabel: the hexadecimal color code for the label should be without the leading #")
			}

			return execError("evalLabel: color code not valid")
		}
	}

	return nil
}

func createLabel(e *Env, labelName *string, label *PadLabel) error {
	err := validateLabelColor(label)
	if err != nil {
		return err
	}

	var labelColor *string
	if label.Color != "" {
		labelColor = &label.Color
	}

	ghLabel := &github.Label{
		Name:        labelName,
		Color:       labelColor,
		Description: &label.Description,
	}

	owner := utils.GetPullRequestOwnerName(e.PullRequest)
	repo := utils.GetPullRequestRepoName(e.PullRequest)

	_, _, err = e.Client.Issues.CreateLabel(e.Ctx, owner, repo, ghLabel)

	return err
}

func checkLabelExists(e *Env, labelName string) (bool, error) {
	owner := utils.GetPullRequestOwnerName(e.PullRequest)
	repo := utils.GetPullRequestRepoName(e.PullRequest)

	_, _, err := e.Client.Issues.GetLabel(e.Ctx, owner, repo, labelName)
	if err != nil {
		if err.(*github.ErrorResponse).Response.StatusCode == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
