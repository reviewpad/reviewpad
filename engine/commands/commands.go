// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands

import (
	"regexp"
)

var Commands = map[*regexp.Regexp]func(matches []string) ([]string, error){
	assignReviewerRegex: AssignReviewer,
}
