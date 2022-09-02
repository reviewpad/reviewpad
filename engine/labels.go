// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"regexp"
	"strings"

	"github.com/google/go-github/v45/github"
)

func validateLabelColor(label *PadLabel) error {
	if label.Color != "" {
		matched, _ := regexp.MatchString(`(?i)^#?([0-9A-F]{6}){1,2}$`, label.Color)
		if !matched {
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
		label.Color = strings.TrimPrefix(label.Color, "#")
		labelColor = &label.Color
	}

	ghLabel := &github.Label{
		Name:        labelName,
		Color:       labelColor,
		Description: &label.Description,
	}

	owner := e.TargetEntity.Owner
	repo := e.TargetEntity.Repo

	_, _, err = e.GithubClient.CreateLabel(e.Ctx, owner, repo, ghLabel)

	return err
}

func checkLabelExists(e *Env, labelName string) (bool, error) {
	owner := e.TargetEntity.Owner
	repo := e.TargetEntity.Repo

	_, _, err := e.GithubClient.GetLabel(e.Ctx, owner, repo, labelName)
	if err != nil {
		if err.(*github.ErrorResponse).Response.StatusCode == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
