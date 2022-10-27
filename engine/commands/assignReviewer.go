// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var assignReviewerRegex = regexp.MustCompile(`^\/reviewpad assign-reviewers\s+((?:[a-zA-Z0-9\-]+[a-zA-Z0-9])(?:,\s*[a-zA-Z0-9\-]+[a-zA-Z0-9])*)(?:\s+(\d+))?(?:\s+(random|round-robin|reviewpad))?$`)

func AssignReviewer(matches []string) ([]string, error) {
	var err error

	if len(matches) < 2 {
		return nil, errors.New("invalid assign reviewer command")
	}

	reviewersList := strings.Split(strings.ReplaceAll(matches[1], " ", ""), ",")

	availableReviewers := `"` + strings.Join(reviewersList, `","`) + `"`

	totalRequiredReviewers := len(reviewersList)

	if len(matches) > 2 && matches[2] != "" {
		totalRequiredReviewers, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, err
		}
	}

	policy := "reviewpad"

	if len(matches) > 3 && matches[3] != "" {
		policy = matches[3]
	}

	action := fmt.Sprintf(`$assignReviewer([%s], %d, %q)`, availableReviewers, totalRequiredReviewers, policy)

	return []string{action}, nil
}
